import 'dart:html';

import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;

void main() {
  runApp(const MyApp());
}

class Case {
  final int CaseID;
  // final String Status; //Case State
  final int Assignee; //Worker UID
  final int CustomerID; 

  Case({
    required this.CaseID,
    // required this.Status,
    required this.Assignee,
    required this.CustomerID,
  });

  factory Case.fromJson(Map<String, dynamic> json) {
    return Case(
      CaseID: json['CaseID'],
      // Status: json['Status'],
      Assignee: json['Assignee']?.isEmpty ?? 0,
      CustomerID: json['CustomerID'],
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
    var d = jsonDecode(response.body); 
    Case c = Case.fromJson(d);  
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
              child: SizedBox(width: 200, child: const Text("Manager"),),
            ),
          ],
        ));
  }
}

GridView _workerGridView(List<Worker> data) {
  return GridView.builder(
    gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
        crossAxisSpacing: 20, maxCrossAxisExtent: 350, mainAxisExtent: 300,  mainAxisSpacing: 20),
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
      const SizedBox(height:20, width: 140, ), 
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
      const Divider(height: 15), 
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

  String getWorkerFace(int workerFaceID) {
    return 'http://localhost:80/worker_faces/FACEID.png'
        .replaceAll("FACEID", workerFaceID.toString());
  }

  void removeDragged(CaseCard toRemove) {
    workerCases
        .removeWhere((e) => e.cardCase.CaseID == toRemove.cardCase.CaseID);
    setState(() {});
  }

  void assignCard(CaseCard c) {
    Map data = Map<String, dynamic>();
    data['workerid'] = widget.worker.WorkerID;
    data['caseid'] = c.cardCase.CaseID;

    http.post(Uri.parse('http://localhost:8001/case/assign'),
        body: json.encode(data));
  }

  @override
  Widget build(BuildContext context) {
    return DragTarget(
      builder: (BuildContext context, List<dynamic> accepted,
          List<dynamic> rejected) {
        return Container(
          color: Colors.black12,
          child: Column(
            children: [
              Row( children: [
                Column(
                  children: [Text('Agent ' + widget.worker.WorkerID.toString() + ': ' + widget.worker.Name),
                  const Divider(height: 15,), 
                  const WorkerSkillTable()],),                 
                Align(alignment: Alignment.topRight,
                  child: Image.network(getWorkerFace(widget.worker.FaceID), width: 100, ),),],),
              const Divider(height:15, thickness:3, indent:10, endIndent: 10, color: Colors.black26,), 
              Column(
                // scrollDirection: Axis.vertical,                
                children: [
                  Column(
                    children: workerCases,
                  ),
                ],
              )
            ],
          ),
        );
      },
      onAccept: (CaseCard draggedCard) {
        assignCard(draggedCard);
        setState(() {
          var newCase = CaseCard(
              cardCase: draggedCard.cardCase,
              onDragComplete: () {
                //Worker lets the API know that the new card is assigned to them
                removeDragged(draggedCard);
              });
          workerCases.add(newCase);
          draggedCard.onDragComplete();
        });
      },
      onWillAccept: (draggedCard) {
        return !workerCases.contains(draggedCard);
      },
    );
  }
}

class WorkerSkillTable extends StatelessWidget {
  const WorkerSkillTable({ Key? key }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 160,
      child: Table(
      columnWidths: const <int, TableColumnWidth> {
        0: FlexColumnWidth(0.25), 
        1: FixedColumnWidth(24), 
        2: FixedColumnWidth(24),
        3: FlexColumnWidth(0.2), 
        4: FixedColumnWidth(24), 
        5: FixedColumnWidth(24), 
        6: FlexColumnWidth(0.3), 
        7: FixedColumnWidth(24), 
        8: FlexColumnWidth(0.15), 
      },
      children: [
        TableRow(children: [
              Container(),
              Icon(Icons.live_help, color: Colors.green,), 
              Icon(Icons.zoom_in, color: Colors.green,), 
              Container(),
              Icon(Icons.admin_panel_settings, color: Colors.green,),
              Icon(Icons.bug_report, color: Colors.green,),
              Container(),
              Icon(Icons.closed_caption, color: Colors.green,),
              Container(),
        ]), 
        TableRow(children: [
              Container(),
              Icon(Icons.construction, color: Colors.green,), 
              Icon(Icons.backup, color: Colors.green,), 
              Container(),
              Icon(Icons.extension, color: Colors.green,),
              Icon(Icons.dashboard, color: Colors.green,),
              Container(),
              Icon(Icons.mood, color: Colors.green,),
              Container(),
        ]), 
      ],
    ))
      ; 
  }
}

class WorkerSkillCard extends StatelessWidget {
  const WorkerSkillCard({ Key? key }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
        // mainAxisSize: MainAxisSize.max,
        children: [
          SizedBox(
            height: 60,
            width: 70,
            child: GridView.count(crossAxisCount: 2,
          // shrinkWrap: true, 
            padding: const EdgeInsets.all(4), 
            children: [
              Icon(Icons.live_help, color: Colors.green,), 
              Icon(Icons.construction, color: Colors.green,), 
              Icon(Icons.zoom_in, color: Colors.green,),
              Icon(Icons.backup, color: Colors.green,),
            ],)),
          SizedBox(
            height: 60,
            width: 70,
            child: GridView.count(crossAxisCount: 2,
          // shrinkWrap: true, 
            padding: const EdgeInsets.all(4), 
            children: [
              Icon(Icons.admin_panel_settings, color: Colors.green,), 
              Icon(Icons.extension, color: Colors.green,), 
              Icon(Icons.bug_report, color: Colors.green,), 
              Icon(Icons.dashboard, color: Colors.green,),              
            ],)), 
          SizedBox(
            height: 60,
            width: 35,
            child: GridView.count(crossAxisCount: 1,
            shrinkWrap: true, 
            // padding: const EdgeInsets.all(4), 
            children: [
              Icon(Icons.closed_caption, color: Colors.green,), 
              Icon(Icons.mood, color: Colors.green,),             
            ],)),  
        
        ],
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

  String getCustomerFace(Case c) {
    return 'http://localhost:80/customer_faces/FACEID_3.png'.replaceAll("FACEID", c.CustomerID.toString()); 
  }

  

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
      onDragCompleted: () => {
        //Remove card from where it came from.
        widget.onDragComplete()}, 
      child: Align( alignment: Alignment.topLeft, 
        child: SizedBox(
          width: 140, height: 60, 
          child: Row(
            children: [
            SizedBox(
              width: 60,
              height: 60,              
              //Customer face Image
              child: Align(alignment: Alignment.center, child: Image.network(getCustomerFace(widget.cardCase),),),
            ),             
            Container(width:4), 
            Container(width: 1, color: Colors.blueGrey), 
            Container(width:4),
            Center(child: Text("Case : " + widget.cardCase.CaseID.toString())), 
        ],) ),)
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
              , onDragComplete: () => {              
                //Remove it from the original card.
                widget.removeFunction(snapshot.data!)} ,); 
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