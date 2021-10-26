import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
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
  
  late Future<List<Worker>> workers;

  void removeCaseCard() {
    
  }

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
            const Align(
              alignment: Alignment.topLeft,
              child: NewCaseColumn(),
            ),
            Expanded(
              child: FutureBuilder<List<Worker>>(
                future: workers,
                builder: (context, AsyncSnapshot snapshot) {
                  if (snapshot.data != null) {
                    return _workerGridView(snapshot.data);
                  } else
                  {
                    return const Center(child: CircularProgressIndicator());
                  }                  
                },
              ),
            ),
            const Align(
              alignment: Alignment.topRight,
              child: Text("Manager"),
            ),
          ],
        ));
  }
}

GridView _workerGridView(List<Worker> data) {
  return GridView.builder(
    gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
        crossAxisSpacing: 20, maxCrossAxisExtent: 400, mainAxisExtent: 300,  mainAxisSpacing: 20),
    itemCount: data.length,
    padding: EdgeInsets.all(10),
    itemBuilder: (BuildContext context, int index) {
      return WorkerCard(worker: data[index],); 
    }
  );
}


class NewCaseColumn extends StatefulWidget {
  const NewCaseColumn({ Key? key }) : super(key: key);

  @override
  _NewCasColumnState createState() => _NewCasColumnState();
}

class _NewCasColumnState extends State<NewCaseColumn> {
  var newCases = <NewCaseCard>[];

  void removeCase(Case c) {
    newCases.removeWhere((element) => element.myCase.CaseID == c.CaseID); 
    setState(() {
      
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(children: [
      ElevatedButton(
        child: const Text('Generate Case'),
        onPressed: () {
          setState(() {
            newCases.add(
              NewCaseCard(futureCase: createCase(), removeFunction: (Case e) => {removeCase(e)}),
            );
          });
        },
      ),
      Column(
        children: newCases,
      ),
    ]);
  }
}


class WorkerCard extends StatefulWidget {
  final Worker worker;

  const WorkerCard({Key? key, required this.worker}) : super(key: key);

  @override
  _WorkerCardState createState() => _WorkerCardState();
}

class _WorkerCardState extends State<WorkerCard> {
  List<CaseCard> workerCases = []; 

  @override
  void initState() {
    super.initState();
  }

  void removeDragged(CaseCard toRemove) {
    workerCases.removeWhere((e) => e.cardCase.CaseID == toRemove.cardCase.CaseID); 
    setState(() {
      
    });
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
                  Column(
                    children: workerCases,
                  ), 
                ],
              )
            );
      }, 
      onAccept: (CaseCard draggedCard) {
        setState(() {
          var newCase = CaseCard(cardCase: draggedCard.cardCase, 
            onDragComplete: () => {removeDragged(draggedCard)}); 
          workerCases.add(newCase); 
          draggedCard.onDragComplete();
        });
      },
    );
  }
}


class CaseCard extends StatefulWidget {
  final Case cardCase; 
  final Function onDragComplete; 

  CaseCard({ Key? key, required this.cardCase, required this.onDragComplete}) : super(key: key);

  @override
  _CaseCardState createState() => _CaseCardState();
}

class _CaseCardState extends State<CaseCard> {

  @override
  Widget build(BuildContext context) {
    return Draggable(
      data: widget,
      feedback: FittedBox(
          fit: BoxFit.contain,
          child: Text("Assign to", 
          textAlign: TextAlign.justify,
          style: TextStyle(
            color: Colors.black12,
            fontWeight: FontWeight.bold,
          ),),
        ), 
      childWhenDragging: SizedBox(height: 100,   child: Text("Assign Case: " + widget.cardCase.CaseID.toString()) ),
      onDragCompleted: () => {widget.onDragComplete()}, 
      child: SizedBox(
        width: 250, 
        height: 100, 
        child: Row(
          children: [
            Align(alignment: Alignment.centerLeft, child: Image.network('http://localhost:80/1_5.png',),), 
            Container(width:4), 
            Container(width: 1, color: Colors.blueGrey), 
            Container(width:4),
            Center(child: Text("Case : " + widget.cardCase.CaseID.toString())), 
        ],) ),
    ); 
  }
}

class NewCaseCard extends StatefulWidget {
  final Future<Case> futureCase;
  final Function removeFunction; 
  late Case myCase; 

  NewCaseCard({ Key? key, required this.futureCase
  , required this.removeFunction}) : super(key: key);

  @override
  State<NewCaseCard> createState() => _NewCaseCardState();
}

class _NewCaseCardState extends State<NewCaseCard> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 30,
      child: FutureBuilder<Case>(
        future: widget.futureCase,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            if (snapshot.hasData) {
              widget.myCase = snapshot.data!; 
              return CaseCard(cardCase: snapshot.data!
              , onDragComplete: () => {widget.removeFunction(snapshot.data!)} ,); 
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