package server

import (
	"fmt"
	"github.com/clstb/phi/ledger/internal/beanacount"
	"github.com/clstb/phi/ledger/internal/config"
	pb "github.com/clstb/phi/proto"
	"io"
	"os"
)

func (s *LedgerServer) DownLoadBeanAccountFile(in *pb.StringMessage, stream pb.BeanAccountService_DownLoadBeanAccountFileServer) error {
	bufferSize := config.DownloadBufferSize
	filePath := beanacount.GetFilePath(in.Value)
	file, err := os.Open(filePath)
	if err != nil {
		s.Logger.Error(err)
		return err
	}
	defer file.Close()

	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		resp := &pb.FileChunkMessage{
			Chunk: buff[:bytesRead],
		}
		err = stream.Send(resp)
		if err != nil {
			s.Logger.Error("error while sending chunk:", err)
			return err
		}
	}
	return nil
}
