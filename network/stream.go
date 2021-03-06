// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package network

import (
	"bufio"
	"errors"
	"reflect"
	"sync"

	"github.com/dappley/go-dappley/network/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/multiformats/go-multiaddr"
	logger "github.com/sirupsen/logrus"
)

const (
	SyncBlock    = "SyncBlock"
	SyncPeerList = "SyncPeerList"
	RequestBlock = "requestBlock"
	BroadcastTx  = "BroadcastTx"
	Unicast      = 0
	Broadcast    = 1
)

var (
	ErrInvalidMessageFormat = errors.New("Message format is invalid")
)

var (
	startBytes = []byte{0x7E, 0x7E}
	endBytes   = []byte{0x7F, 0x7F, 0}
)

type Stream struct {
	node       *Node
	peerID     peer.ID
	remoteAddr multiaddr.Multiaddr
	stream     net.Stream
	dataCh     chan []byte
	quitRdCh   chan bool
	quitWrCh   chan bool
}

func NewStream(s net.Stream, node *Node) *Stream {
	return &Stream{node,
		s.Conn().RemotePeer(),
		s.Conn().RemoteMultiaddr(),
		s,
		make(chan []byte, 5), //TODO: Redefine the size of the channel
		make(chan bool, 1),   //two channels to stop
		make(chan bool, 1),
	}
}

func (s *Stream) Start() {
	rw := bufio.NewReadWriter(bufio.NewReader(s.stream), bufio.NewWriter(s.stream))
	s.startLoop(rw)
}

func (s *Stream) StopStream() {
	logger.Debug("Stream Terminated! Peer Addr:", s.remoteAddr)
	s.quitRdCh <- true
	s.quitWrCh <- true
	s.stream.Close()
	delete(s.node.streams, s.peerID)
	s.node.peerList.DeletePeer(&Peer{s.peerID, s.remoteAddr})
}

func (s *Stream) Send(data []byte) {
	s.dataCh <- data
}

func (s *Stream) startLoop(rw *bufio.ReadWriter) {
	go s.readLoop(rw)
	go s.writeLoop(rw)
}

func readMsg(rw *bufio.ReadWriter) ([]byte, error) {
	var bytes []byte
	for {
		b, err := rw.ReadByte()

		if err != nil {
			return bytes, err
		}
		bytes = append(bytes, b)
		if containEndingBytes(bytes) {
			return bytes, nil
		}
	}
}

func (s *Stream) read(rw *bufio.ReadWriter) {
	//read stream with delimiter
	bytes, err := readMsg(rw)

	if err != nil {
		logger.Warn(err)
	}

	//TODO: How to verify the integrity of the received message
	if len(bytes) > 1 {
		s.parseData(bytes)
	} else {
		logger.Debug("Read less than 1 byte. Stop Reading...")
		//stop the stream
		s.StopStream()
	}

}

func (s *Stream) readLoop(rw *bufio.ReadWriter) {
	for {
		select {
		case <-s.quitRdCh:
			logger.Debug("Stream ReadLoop Terminated!")
			return
		default:
			s.read(rw)
		}
	}
}

func encodeMessage(data []byte) []byte {
	ret := append(startBytes, data...)
	ret = append(ret, endBytes...)
	return ret
}

func decodeMessage(data []byte) ([]byte, error) {
	if !(containStartingBytes(data) && containEndingBytes(data)) {
		return nil, ErrInvalidMessageFormat
	}
	return data[2 : len(data)-3], nil
}

func containStartingBytes(data []byte) bool {
	if len(data) < len(startBytes) {
		return false
	}
	return reflect.DeepEqual(data[0:len(startBytes)], startBytes)
}

func containEndingBytes(data []byte) bool {
	if len(data) < len(endBytes) {
		return false
	}
	return reflect.DeepEqual(data[(len(data)-len(endBytes)):], endBytes)
}

func (s *Stream) writeLoop(rw *bufio.ReadWriter) error {
	var mutex = &sync.Mutex{}
	for {
		select {
		case data := <-s.dataCh:
			mutex.Lock()
			//attach a delimiter byte of 0x00 to the end of the message
			rw.WriteString(string(encodeMessage(data)))
			rw.Flush()
			mutex.Unlock()
		case <-s.quitWrCh:
			logger.Debug("Stream Write Terminated!")
			return nil
		}
	}
	return nil
}

//should parse and relay
func (s *Stream) parseData(data []byte) {

	dataDecoded, err := decodeMessage(data)
	if err != nil {
		logger.Warn(err)
		return
	}

	dmpb := &networkpb.Dapmsg{}
	//unmarshal byte to proto
	if err := proto.Unmarshal(dataDecoded, dmpb); err != nil {
		logger.Info(err)
	}

	dm := &DapMsg{}
	dm.FromProto(dmpb)

	switch dm.GetCmd() {
	case SyncBlock:
		logger.Debug("Stream: Received ", SyncBlock, " command from:", dm.key)
		s.node.syncBlockHandler(dm, s.peerID)
	case SyncPeerList:
		logger.Debug("Stream: Received ", SyncPeerList, " command from:", s.remoteAddr)
		s.node.addMultiPeers(dm.GetData())
	case RequestBlock:
		logger.Debug("Stream: Received ", RequestBlock, " command from:", s.remoteAddr)
		s.node.sendRequestedBlock(dm.GetData(), s.peerID)
	case BroadcastTx:
		logger.Debug("Stream: Received ", BroadcastTx, " command from:", s.remoteAddr)
		s.node.addTxToPool(dm.GetData())
	default:
		logger.Debug("Received invalid command from:", s.remoteAddr)
	}

}
