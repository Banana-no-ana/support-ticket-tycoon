import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;

void main() {
  runApp(const MyApp());
}

class Case {
  final int CaseID;
  final String State; //Case State
  final int Assignee; //Worker UID

  Case({
    required this.CaseID,
    required this.State,
    required this.Assignee,
  });

  factory Case.fromJson(Map<String, dynamic> json) {
    return Case(
      CaseID: json['CaseID'],
      State: json['State'],
      Assignee: json['Assignee'],
    );
  }
}

Future<List<Worker>> getWorkers() async {
  final response = await http.get(
    Uri.parse('http://localhost:8001/worker/list'),
  );

  late List<Worker> workers = [];

  if (response.statusCode == 200) {
    jsonDecode(response.body).forEach((element) {
      workers.add(Worker.fromJson(element));
    });

    return workers;
  } else {
    // If the server did not return a 200 OK response,
    // then throw an exception.
    throw Exception('Failed to get Case');
  }
}

Future<Case> createCase() async {
  final response = await http.get(
    Uri.parse('http://localhost:8001/case/create'),
  );

  if (response.statusCode == 200) {
    // If the server did return a 200 OK response,
    // then parse the JSON.
    return Case.fromJson(jsonDecode(response.body));
  } else {
    // If the server did not return a 200 OK response,
    // then throw an exception.
    throw Exception('Failed to get Case');
  }
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Scenario 0',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const MyHomePage(title: 'Scenario 0'),
    );
  }
}


class Worker {
  final int WorkerID; //Assigned worker ID
  final String Name; //Generated from a list of names
  final int FaceID; //Icon for worker face

  Worker({
    required this.WorkerID,
    required this.Name,
    required this.FaceID,
  });

  factory Worker.fromJson(Map<String, dynamic> json) {
    return Worker(
      WorkerID: json['WorkerID'],
      Name: json['Name'],
      FaceID: json['FaceID'],
    );
  }
}


class MyHomePage extends StatefulWidget {
  const MyHomePage({Key? key, required this.title}) : super(key: key);
  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  var newCases = <NewCaseCard2>[];
  late Future<List<Worker>> workers;

  @override
  void initState() {
    super.initState();
    workers = getWorkers();    
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text(widget.title),
        ),
        body: Row(
          children: [
            Align(
              alignment: Alignment.topLeft,
              child: Column(children: [
                ElevatedButton(
                  child: const Text('Generate Case'),
                  onPressed: () {
                    setState(() {
                      newCases.add(NewCaseCard2(futureCase: createCase()));
                    });
                  },
                ),
                Column(
                  children: newCases,
                ),
              ]),
            ),
            Expanded(
              child: FutureBuilder<List<Worker>>(
                future: workers,
                builder: (context, AsyncSnapshot snapshot) {
                  if (snapshot.data != null) {
                    return _workerGridView(snapshot.data);
                  } else
                  {
                    return Center(child: CircularProgressIndicator());
                  }                  
                },
              ),
            ),
            Align(
              alignment: Alignment.topRight,
              child: Text("Manager"),
            ),
          ],
        ));
  }
}

GridView _workerGridView(List<Worker> data) {
  return GridView.builder(
    gridDelegate: SliverGridDelegateWithMaxCrossAxisExtent(
        crossAxisSpacing: 20, maxCrossAxisExtent: 400, mainAxisExtent: 300,  mainAxisSpacing: 20),
    itemCount: data.length,
    padding: EdgeInsets.all(10),
    itemBuilder: (BuildContext context, int index) {
      return WorkerCard(worker: data[index],); 
    }
  );
}


class WorkerCard extends StatefulWidget {
  final Worker worker;

  const WorkerCard({Key? key, required this.worker}) : super(key: key);

  @override
  _WorkerCardState createState() => _WorkerCardState();
}

class _WorkerCardState extends State<WorkerCard> {
  List<Case> cases = [];

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
      return DragTarget(
        builder: (BuildContext context, List<dynamic> accepted,
          List<dynamic> rejected) {
        return Container(
              width: 30,
              height: 30, 
              color: Colors.amber,
              child: Column(
                children: [
                  Text('#' + widget.worker.WorkerID.toString() + ' ' + widget.worker.Name),
                ],
              )
            );
      }, 
      onAccept: (Case draggedCase) {
        setState(() {
          cases.add(draggedCase);  
        });
        
      },
    );
  }
}


class CaseCard extends StatefulWidget {
  final Case cardCase; 

  CaseCard({ Key? key, required this.cardCase}) : super(key: key);

  @override
  _CaseCardState createState() => _CaseCardState();
}

class _CaseCardState extends State<CaseCard> {

  @override
  Widget build(BuildContext context) {
    return Draggable(
      data: widget.cardCase,
      feedback: SizedBox(height: 40, width: 100, child: Text("Assign to"), ),
      childWhenDragging: SizedBox(height: 40,  child: Text("Assign Case: " + widget.cardCase.CaseID.toString()) ),
      child: SizedBox(
        width: 100, 
        height: 40, 
        child: Column(children: [
          Center(child: Text("Case : " + widget.cardCase.CaseID.toString())), 
        ],) ),
    ); 
  }
}

class NewCaseCard2 extends StatelessWidget {
  final Future<Case> futureCase;

  const NewCaseCard2({ Key? key, required this.futureCase }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 30,
      child: FutureBuilder<Case>(
        future: futureCase,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            if (snapshot.hasData) {
              return CaseCard(cardCase: snapshot.data!); 
             }
          }
          else if (snapshot.connectionState == ConnectionState.waiting) {
            return Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[Text("Case being created")],
              );
          }
          return const CircularProgressIndicator();
        },
      ), 

    );
  }
}