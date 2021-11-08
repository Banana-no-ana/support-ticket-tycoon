import pytest
import requests
import subprocess
import time
import json

apiaddr = "http://localhost:8001/"

class TestStartup: 
    def test_alive(self): 
        r = requests.get(apiaddr + "healthz")
        assert r.text == "ok" 


class TestCases: 
    def test_create(self):
        r = requests.get(apiaddr + "case/create")
        assert r.status_code == 200
        assert r.json()['CaseID'] >= 0

    def test_list(self):
        r = requests.get(apiaddr + "case/list")
        assert r.status_code == 200
        assert len(r.json()) > 0


class TestStartWorker:
    def test_startworker(self):
        r = requests.get(apiaddr + "worker/create")
        assert r.text.startswith("Created worker")

    def test_addworker(self):
        subprocess.Popen(["go", "run", "../worker.go", "--rpc_port=:10202", "--http_port=:9202", "--worker_id=202"] )
        time.sleep(0.5)

        worker = {'WorkerID': 202, 'Name': "Test Worker 202", 'FaceID': 1}
        with requests.post(apiaddr + "worker/add", json.dumps(worker)) as r: 
            assert r.status_code == 200
            assert r.text.__contains__("202")
    
    def test_addworker_assign(self): 
        r = requests.get(apiaddr + "case/create")
        caseid = r.json()['CaseID']

        data = {'CaseID':caseid , 'WorkerID':202}
        with requests.post(apiaddr+ "case/assign", json.dumps(data)) as r: 
            assert r.status_code == 200
            assert r.text.__contains__("successful")

        with requests.get("http://localhost:9202/case/list") as r: 
            assert r.text.__contains__(str(caseid))



class TestScenario:
    def test_loadscenario(self): 
        r = requests.get(apiaddr + "scenario/load/0") #Scenario 0 should have 4 workers
        assert r.text.__contains__("Created 4 Worker")

    def test_workerSkill(self): 
        r = requests.get("http://localhost:9001/skill/list")
        assert r.status_code == 200
        assert r.text.__contains__("Build")

    def testAssign(self):
        r = requests.get(apiaddr + "case/create")
        caseid = r.json()['CaseID']

        data = {'CaseID':caseid , 'WorkerID':1}
        with requests.post(apiaddr+ "case/assign", json.dumps(data)) as r: 
            assert r.status_code == 200
            assert r.text.__contains__("successful")

        with requests.get("http://localhost:9001/case/list") as r: 
            assert r.text.__contains__(str(caseid))

       
    def test_UnloadScenario(self):
        with requests.get(apiaddr + "scenario/unload") as r: 
            assert r.status_code == 200
            assert r.text.__contains__("Removed")        
    
        time.sleep(1)
        with pytest.raises(Exception) as e:
            requests.get("http://localhost:9001/healthz")
