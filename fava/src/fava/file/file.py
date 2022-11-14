import logging

import grpc
#from .shared_pb2 import StringMessage
from .ledger_pb2_grpc import BeanAccountServiceStub
import os

LEDGER_URI = (
    os.environ["LEDGER_URI"] if "LEDGER_URI" in os.environ else "localhost:8082"
)
DATA_DIR_PATH = (
    os.environ["DATA_DIR_PATH"] if "DATA_DIR_PATH" in os.environ else "../.data"
)

log = logging.getLogger(__name__)


def provision_file(username: str):
    with grpc.insecure_channel(LEDGER_URI) as channel:
        stub = BeanAccountServiceStub(channel)
        log.info(f"-------> DownLoadBeanAccountFile({username})")
        message = StringMessage(value=username)
        rel_path = f"{DATA_DIR_PATH}/{username}.beancount"
        with open(rel_path, "wb+") as bfile:
            for res in stub.DownLoadBeanAccountFile(message):
                bfile.write(res.Chunk)
        log.info("<------- OK")
        return os.path.abspath(rel_path)
