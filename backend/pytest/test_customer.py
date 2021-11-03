import pytest
import requests
import subprocess
import time
import logging
import json

from requests import api
from requests.sessions import HTTPAdapter


workeraddr = "http://localhost:9201/"
adpiaddr = "http://localhost:8001/"
custaddr = "http://localhost:8005/"

class TestCustomer:
    def test_start_test(self): 
        try:
            requests.get(custaddr + "healthz")
        except:
            subprocess.Popen(["go", "run", "../customer.go"] )
            time.sleep(1)
        
        with requests.get(custaddr + "healthz") as r:
            assert r.text == "ok"   

    def test_list_cases(self):
        requests.get(adpiaddr + "case/create")        

        with requests.get(custaddr + "case/list") as r:
            assert r.status_code == 200
            assert len(r.json()) > 0
            assert r.text.__contains__("CustomerID")

    def test_manual_register(self):
        with requests.get(custaddr + "case/register/100") as r:
            assert r.status_code == 200
        
        with requests.get(custaddr + "case/list") as r:
            assert r.text.__contains__("100")
            

    def test_end_test(self): 
        requests.get(custaddr + "kill")