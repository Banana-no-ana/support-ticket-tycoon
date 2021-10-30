import pytest
import requests
import subprocess

workeraddr = "http://localhost:9201/"

class TestWorker:    
    def setup_method(self, test_method):
        # os.system("go run ../worker.go --rpc_port=:10200 --http_port=:9200 --worker_id=200 ") 
        subprocess.Popen(["go", "run", "../worker.go", "--rpc_port=:10201", "--http_port=:9201", "--worker_id=201"] )

    def teardown_method(self, test_method):
        r = requests.get(workeraddr + "kill")

    def test_startup(self):         
        with requests.get(workeraddr + "healthz") as r:
            assert r.text == "ok"   

    def test_assign(self):
        with requests.get(workeraddr + "assign/100") as r: 
            assert r.status_code == 200
            assert r.text.startswith("case accepted")

    def test_listcases(self):
        with requests.get(workeraddr + "case/list") as r:
            assert r.status_code == 200
            assert r.text.__contains__("100")
