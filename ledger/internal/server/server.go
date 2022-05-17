package server

import (
	"github.com/clstb/phi/pkg"
	pb "github.com/clstb/phi/proto"
	"go.uber.org/zap"
)

type LedgerServer struct {
	pb.BeanAccountServiceServer
	Logger      *zap.SugaredLogger
	TinkGwUri   string
	DataDirPath string
}

func NewServer(tinkGwUri string, dataDirPath string) *LedgerServer {
	sugar := pkg.CreateLogger()
	s := &LedgerServer{
		Logger:      sugar.Named("Ledger"),
		TinkGwUri:   tinkGwUri,
		DataDirPath: dataDirPath,
	}
	return s
}
