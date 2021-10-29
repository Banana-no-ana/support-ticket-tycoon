import pytest
import requests
import subprocess

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
       

class TestAssignCaseToWorker:
    def testAssign(self):
        data = {'caseid':1 , 'workerid':1}
        r = requests.post(apiaddr+ "case/assign", data)
        assert r.status_code == 200
        assert r.text.__contains__("sucessful")

        # r2 = requests.get("http://localhost/") //Need to also check that the worker actually has it


