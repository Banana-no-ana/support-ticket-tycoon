import pytest
import requests
import subprocess

workeraddr = "http://localhost:9200/"

class TestStartup:    

    def setup_method(self, test_method):
        # os.system("go run ../worker.go --rpc_port=:10200 --http_port=:9200 --worker_id=200 ") 
        subprocess.Popen(["go", "run", "../worker.go", "--rpc_port=:10200", "--http_port=:9200", "--worker_id=200"] )

    def teardown_method(self, test_method):
        r = requests.get(workeraddr + "kill")

    def test_startup(self):         
        r = requests.get(workeraddr + "healthz")
        assert r.text == "ok" 

    def test_assign(self):
        r = requests.get(workeraddr + "assign/1")
        assert r.status_code == 200
        assert r.text.startswith("case accepted")
