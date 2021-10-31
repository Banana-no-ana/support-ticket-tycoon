import pytest
import requests
import subprocess
import time

apiaddr = "http://localhost:8001/"

class TestStartup: 
    def test_alive(self): 
        r = requests.get(apiaddr + "healthz")
        assert r.text == "ok" 


class TestCases: 
    def test_create(self):
        r = requests.get(apiaddr + "case/create")
        assert r.status_code == 200
        assert r.json()['State'] == "New"

    def test_list(self):
        r = requests.get(apiaddr + "case/list")
        assert r.status_code == 200
        assert len(r.json()) > 0


class TestStartWorker:
    def test_startworker(self):
        r = requests.get(apiaddr + "worker/create")
        assert r.text.startswith("Created worker")


class TestScenario:
    def test_loadscenario(self): 
        r = requests.get(apiaddr + "scenario/load/0") #Scenario 0 should have 4 workers
        assert r.text.__contains__("Created 4 Worker")

        time.sleep(2)

    # def test_workerSkill(self): 
    #     r = requests.get("http://localhost:9001/skill/list")
    #     assert r.status_code == 200
    #     assert r.text.__contains__("Build")

    def testAssign(self):
        data = {'caseid':1 , 'workerid':1}
        with requests.post(apiaddr+ "case/assign", data) as r: 
            assert r.status_code == 200
            assert r.text.__contains__("successful")

    def testAssign2(self):
        data = {'caseid':100 , 'workerid':2}
        with requests.post(apiaddr+ "case/assign", data) as r: 
            assert r.status_code == 200
        
        with requests.get("http://localhost:9002/case/list") as r: 
            assert r.text.__contains__("100")
       
    def test_UnloadScenario(self):
        with requests.get(apiaddr + "scenario/unload") as r: 
            assert r.status_code == 200
            assert r.text.__contains__("Removed")        
        
        time.sleep(1)
        with pytest.raises(Exception) as e:
            requests.get("http://localhost:9001/healthz")
