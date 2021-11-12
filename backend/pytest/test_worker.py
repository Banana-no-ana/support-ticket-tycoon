import pytest
import requests
import subprocess
import time
import logging
import json

from requests import api


workeraddr = "http://localhost:9201/"
adpiaddr = "http://localhost:8001/"

class TestWorker:    
    def test_start_test(self): 
        subprocess.Popen(["go", "run", "../worker.go", "--rpc_port=:10201", "--http_port=:9201", "--worker_id=201"] )
        time.sleep(1)
        with requests.get(workeraddr + "healthz") as r:
            assert r.text == "ok"   

    def test_assign(self):
        with requests.get(workeraddr + "case/assign/100") as r:
            assert r.status_code == 200
            assert r.text.startswith("case accepted")

    def test_listcases(self, tolist='100'):
        with requests.get(workeraddr + "case/list") as r: 
            assert r.status_code == 200
            assert r.text.__contains__(tolist)

    def test_unassign(self): 
        with requests.get(workeraddr + "case/unassign/100") as r:
            assert r.status_code == 200
            assert r.text.startswith("200")

    def test_unassign_fail(self): 
        with requests.get(workeraddr + "case/unassign/100") as r:
            assert r.status_code == 200
            assert r.text.startswith("404")

    def test_endtest(self):
        r = requests.get(workeraddr + "kill")