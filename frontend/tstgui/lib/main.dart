import 'dart:async';
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
  final int CustomerSentiment;
  final String Status; 
  final List<CaseStage> CaseStages; 

  Case({
    required this.CaseID,
    // required this.Status,
    required this.Assignee,
    required this.CustomerID,
    required this.CustomerSentiment,
    required this.Status, 
    required this.CaseStages,
  });

  factory Case.fromJson(Map<String, dynamic> json) {
    List<CaseStage> stages = []; 
    if (json.containsKey('CaseStages')) {
      Iterable stageList = json['CaseStages']; 
      // CaseStage firStage = CaseStage.fromJson(stageList.first); 
      stages = List<CaseStage>.from(stageList.map((e) => CaseStage.fromJson(e)));   
    }   

    return Case(
      CaseID: json['CaseID'],
      // Status: json['Status'],
      Assignee: json['Assignee'] ?? 0,
      CustomerID: json['CustomerID'],
      CustomerSentiment: json['CustomerSentiment'] ?? 3, 
      Status: json['Status'] ?? "New",
      CaseStages: stages, 
    );
  }
}

class CaseStage {
  final int StageID; 
  final String StageType; 
  final String StageStatus;
  final int CompletedWork; 
  final int TotalWork; 

  CaseStage({
    required this.StageID, 
    required this.StageType, 
    required this.StageStatus, 
    required this.CompletedWork, 
    required this.TotalWork, 
  }); 

  factory CaseStage.fromJson(Map<String, dynamic> json) {
    return CaseStage(
      StageID: json['StageID']?? 0, 
      StageType: json['Type']?? 'Undefined',
      // Status: json['Status'],
      StageStatus: json['Status'] ?? '',
      CompletedWork: json['Completedwork'] ?? 0,
      TotalWork: json['Totalwork'] ?? 0, 
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
  final WorkerSkill Skills; 

  Worker({
    required this.WorkerID,
    required this.Name,
    required this.FaceID,
    required this.Skills, 
  });

  factory Worker.fromJson(Map<String, dynamic> json) {
    return Worker(
      WorkerID: json['WorkerID'],
      Name: json['Name'],
      FaceID: json['FaceID'],
      Skills: WorkerSkill.fromJson(json['Skills'],)
    );
  }
}

class WorkerSkill {
  int Troubleshoot = 1;
  int Build = 2;
  int Research = 3;
  int WebTech = 4;
  
  int Admin = 5;
  int Usage = 6; 
  int Architecture = 7; 
  int Ecosystem = 8; 

  int Explain = 9;
  int Empathy = 10;

  WorkerSkill({
    required this.Troubleshoot,
    required this.Build,
    required this.Research,
    required this.WebTech,
    required this.Admin,
    required this.Usage,
    required this.Architecture,
    required this.Ecosystem,
    required this.Explain,
    required this.Empathy,
  });

  factory WorkerSkill.fromJson(Map<String, dynamic> json) {
    return WorkerSkill(
      Troubleshoot: json['Troubleshoot'] ?? 0,
      Build: json['Build'] ?? 0 ,
      Research: json['Research'] ?? 0,
      WebTech: json['WebTech'] ?? 0,
      Admin: json['Admin'] ?? 0,
      Usage: json['Usage'] ?? 0,
      Architecture: json['Architecture'] ?? 0,
      Ecosystem: json['Ecosystem'] ?? 0,
      Explain: json['Explain'] ?? 0,
      Empathy: json['Empathy'] ?? 0,    
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
        crossAxisSpacing: 20, maxCrossAxisExtent: 350, mainAxisExtent: 350,  mainAxisSpacing: 20),
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

class _WorkerCardState extends State<WorkerCard> with TickerProviderStateMixin {
  List<CaseCard2> workerCases = [];
  List<CompletedCaseCard> completedCases = [];

  @override
  void initState() {
    super.initState();
  }

  String getWorkerFace(int workerFaceID) {
    return 'http://localhost:80/worker_faces/FACEID.png'
        .replaceAll("FACEID", workerFaceID.toString());
  }

  void removeDragged(CaseContainer toRemove) {
    workerCases
        .removeWhere((e) => e.cardCase.CaseID == toRemove.cardCase.CaseID);
    setState(() {});
  }

  void moveToCompleted(CaseContainer toMove) {
    completedCases.add(
      CompletedCaseCard(cardCase: toMove.cardCase,) 
    ); 
    workerCases
        .removeWhere((e) => e.cardCase.CaseID == toMove.cardCase.CaseID);
    setState(() {});
  }

  void assignCard(CaseContainer c) {
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
                  WorkerSkillTable(skills: widget.worker.Skills)],),                 
                Align(alignment: Alignment.topRight,
                  child: Image.network(getWorkerFace(widget.worker.FaceID), width: 100, ),),],),
                const Divider(height:15, thickness:3, indent:10, endIndent: 10, color: Colors.black26,), 
                SizedBox(height: 150, 
                  child: ListView.builder(
                    itemCount: workerCases.length,
                    itemBuilder: (BuildContext context, int index) {
                      return Container(child: workerCases[index]); 
                  }), ), 
                const Divider(height: 5), 
                SizedBox(height: 80, 
                  child: ListView.builder(
                    itemCount: completedCases.length,
                    itemBuilder: (BuildContext context, int index) {
                      return Container(child: completedCases[index]); 
                  }), ), 
            ],
          ),
        );
      },
      onAccept: (CaseContainer draggedCard) {
        assignCard(draggedCard);
        setState(() {
          var newCase = CaseCard2(
              cardCase: draggedCard.cardCase,
              initialIncrupProgressValue: draggedCard.incrupController.value, 
              onCaseComplete: () => {moveToCompleted(draggedCard)},
              onDragComplete: () {
                //the target draggable with call this function so to remove it from the existing worker
                removeDragged(draggedCard);
              }, );
          // newCase.incrupController.forward();  
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
  final WorkerSkill skills; 
  const WorkerSkillTable({ Key? key , required this.skills }) : super(key: key);

  Color gs(int skillLevel) {
    switch (skillLevel) {
      case 0: return Colors.black12; 
      case 1: return Colors.white; 
      case 2: return Colors.yellow; 
      case 3: return Colors.orange; 
      case 4: return Colors.blue; 
      case 5: return Colors.black; 
    }
    return Colors.black12; 
  }

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
              Icon(Icons.live_help, color: gs(skills.Troubleshoot),), 
              Icon(Icons.zoom_in, color: gs(skills.Research),), 
              Container(),
              Icon(Icons.admin_panel_settings, color: gs(skills.Admin),),
              Icon(Icons.bug_report, color: gs(skills.Architecture),),
              Container(),
              Icon(Icons.closed_caption, color: gs(skills.Explain),),
              Container(),
        ]), 
        TableRow(children: [
              Container(),
              Icon(Icons.construction, color:  gs(skills.Build),), 
              Icon(Icons.backup, color:  gs(skills.WebTech),), 
              Container(),
              Icon(Icons.extension, color:  gs(skills.Usage),),
              Icon(Icons.dashboard, color:  gs(skills.Ecosystem),),
              Container(),
              Icon(Icons.mood, color: gs(skills.Empathy),),
              Container(),
        ]), 
      ],
    ))
      ; 
  }
}

class CaseProgressCard extends StatelessWidget {
  final List<CaseStage> stages; 
  const CaseProgressCard({ Key? key, required this.stages }) : 
    super(key: key);

  Color getIconColor(int difficulty) {
    switch (difficulty) {
      case 0: return Colors.black12; 
      case 1: return Colors.white; 
      case 2: return Colors.yellow; 
      case 3: return Colors.orange; 
      case 4: return Colors.blue; 
      case 5: return Colors.black; 
      default: return Colors.black12; 
    }
  }

  Stack buildSkillIcon(CaseStage sc) {
    Icon ic;
    double s = 32; 
    switch (sc.StageType) {
      case "Troubleshoot":
        ic = Icon(Icons.live_help, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Build":
        ic = Icon(Icons.construction, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Research":
        ic = Icon(Icons.zoom_in, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "WebTech":
        ic = Icon(Icons.backup, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Admin":
        ic = Icon(Icons.admin_panel_settings, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Usage":
        ic = Icon(Icons.extension, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Architecture":
        ic = Icon(Icons.bug_report, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Ecosystem":
        ic = Icon(Icons.dashboard, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Explain":
        ic = Icon(Icons.closed_caption, color: getIconColor(sc.TotalWork), size: s); 
        break;
      case "Empathy":
        ic = Icon(Icons.mood, color: getIconColor(sc.TotalWork), size: s); 
        break;
      default:
        ic = const Icon(Icons.error); 
    }

    return Stack( children: [
      Center(child: SizedBox(width: 40, height: 40, child: ic), ), 
      Center(child: SizedBox(width: 41, height: 42, child: CircularProgressIndicator(value: sc.CompletedWork / sc.TotalWork), ),)
      ]);      
  }

  @override
  Widget build(BuildContext context) {
    return Center(child: ListView.separated(
      scrollDirection: Axis.horizontal,
      physics: const NeverScrollableScrollPhysics(),
      itemBuilder: (BuildContext context, int index) {
         return buildSkillIcon(stages[index]); 
      }, 
      separatorBuilder:  (BuildContext context, int index) => Container(width: 7,), 
      itemCount: stages.length),);
  }
}

//Define interface for new casecard and old case cards
abstract class CaseContainer {
  final Case cardCase; 
  final Function onDragComplete;
  final Function onCaseComplete;
  late AnimationController incrupController;
  final double initialIncrupProgressValue; 

  CaseContainer({ Key? key, required this.cardCase, 
    required this.onDragComplete,
    required this.initialIncrupProgressValue, 
    required this.onCaseComplete}); 
}

class CompletedCaseCard extends StatelessWidget {
  final Case cardCase; 

  const CompletedCaseCard({ Key? key, 
    required this.cardCase, }) : super(key: key);
  
  @override
  Widget build(BuildContext context) {
    return Container(
      child: Text("Completed Case"),
    );
  }
}

//New case card because dynamically figuring out sizing new cases and existing cases is too hard :(
class CaseCard2 extends StatefulWidget implements CaseContainer{
  final Case cardCase; 
  final Function onDragComplete;
  final Function onCaseComplete;
  late AnimationController incrupController;
  final double initialIncrupProgressValue; 

  CaseCard2({ Key? key, required this.cardCase, 
    required this.onCaseComplete,   
    required this.onDragComplete,required this.initialIncrupProgressValue, }) : super(key: key);

  @override
  _CaseCard2State createState() => _CaseCard2State();
}

class _CaseCard2State extends State<CaseCard2> with TickerProviderStateMixin {
  late Timer caseUpdateTimer; 
  late List<CaseStage> caseStages; 

  String getCustomerFace(Case c) {
    return 'http://localhost:80/customer_faces/FACEID_CUSTOMER_SENTIMENT.png'.
    replaceAll("FACEID", c.CustomerID.toString()).
    replaceAll("CUSTOMER_SENTIMENT", c.CustomerSentiment.toString()); 
  }


  Future<void> getCaseUpdate(Case supCase) async {
    final response = await http.get(
      Uri.parse('http://localhost:8001/case/get/CASEID'
        .replaceAll('CASEID', supCase.CaseID.toString())),
    );

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response,
      // then parse the JSON.
      var c = Case.fromJson(jsonDecode(response.body)); 
      // caseStages = c.CaseStages; 
      caseStages = c.CaseStages; 

      if (c.Status == "Closed") {
        widget.onCaseComplete(); 
        caseUpdateTimer.cancel(); 
      }

      return;
    } else {
      // If the server did not return a 200 OK response,
      // then throw an exception.
      throw Exception('Failed to get Case');
    }
  }

  @override
  void initState() {
    caseStages = []; 
    widget.incrupController = AnimationController(vsync: this, 
    duration: const Duration(seconds:8))..addListener(() {
        setState(() {});
      });
    caseUpdateTimer = Timer.periodic(Duration(milliseconds: 600), (Timer t) {
      //First make the network call to get case updates. Then setstate       
      setState(() {
        unawaited(getCaseUpdate(widget.cardCase)); 
      });
     }); 
    widget.incrupController.value = widget.initialIncrupProgressValue; 
    widget.incrupController.forward();
    super.initState();
  }

  @override
  void dispose() {
    widget.incrupController.dispose(); 
    caseUpdateTimer.cancel(); 
    super.dispose();
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
          height: 70, 
          child: Row(
            children: [
            Container(width: 3), 
            SizedBox(
              width: 60,
              height: 70,
              //Customer face Image
              child: Column(
                children: [
                  SizedBox( width: 60,
                    height: 60,  child: Align(alignment: Alignment.center, 
                    child: Image.network(getCustomerFace(widget.cardCase),),),), 
                  LinearProgressIndicator(value: widget.incrupController.value,), 
                ],)
            ),
            Container(width:4), 
            Container(width: 1, color: Colors.blueGrey), 
            Container(width:4),
            SizedBox(
              width: 200, 
              child: Center(
                child: Column(children: [
                  Text("Case : " + widget.cardCase.CaseID.toString()),
                  const Divider(height: 7), 
                  Center(child: SizedBox(height: 40, child: 
                        CaseProgressCard(stages: caseStages),),), 
                    // SizedBox(height: 24, width: 200, child: 
              //     Column(children: [                
              //       InkWell(child: Text("Get Case Update"), 
              //         onTap: () {
              //           unawaited(getCaseUpdate(widget.cardCase)); 
              //           setState(() => {});
              //         },),                   
              // ],),               
              // ), 
            ] ,), 
              ),)             
        ],) ),)
    ); 
  }
}


class CaseCard extends StatefulWidget implements CaseContainer {
  final Case cardCase; 
  final Function onDragComplete;
  final Function onCaseComplete;
  late AnimationController incrupController;
  final double initialIncrupProgressValue; 


  CaseCard({ Key? key, required this.cardCase, 
    required this.onCaseComplete, 
    required this.onDragComplete,required this.initialIncrupProgressValue, }) : super(key: key);

  @override
  _CaseCardState createState() => _CaseCardState();
}

class _CaseCardState extends State<CaseCard> with TickerProviderStateMixin{
  String getCustomerFace(Case c) {
    return 'http://localhost:80/customer_faces/FACEID_CUSTOMER_SENTIMENT.png'.
    replaceAll("FACEID", c.CustomerID.toString()).
    replaceAll("CUSTOMER_SENTIMENT", c.CustomerSentiment.toString()); 
  }

  @override
  void initState() {
    widget.incrupController = AnimationController(vsync: this, 
    duration: const Duration(seconds:8))..addListener(() {
        setState(() {});
      });
    widget.incrupController.value = widget.initialIncrupProgressValue; 
    widget.incrupController.forward();
    super.initState();
  }

  @override
  void dispose() {
    widget.incrupController.dispose(); 
    super.dispose();
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
          width: 140, height: 100, 
          child: Row(
            children: [
            SizedBox(
              width: 60,
              height: 100,              
              //Customer face Image
              child: Align(alignment: Alignment.center, child: Image.network(getCustomerFace(widget.cardCase),),),
            ),             
            Container(width:4), 
            Container(width: 1, color: Colors.blueGrey), 
            Container(width:4),
            SizedBox(
              width: 70,
              child: Center(child: Column(children: [
              Text("Help!",), 
              LinearProgressIndicator(value: widget.incrupController.value,), 
            ] ,), 
              ),)             
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

class _NewCaseCardState extends State<NewCaseCard> with TickerProviderStateMixin {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 60,
      child: FutureBuilder<Case>(
        future: widget.futureCase,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            if (snapshot.hasData) {
              widget.myCase = snapshot.data!; 
              var card = CaseCard(cardCase: snapshot.data!, 
                initialIncrupProgressValue: 0.0,
                onCaseComplete: () => {},
                onDragComplete: () => {              
                //Remove it from the original card.
                widget.removeFunction(snapshot.data!)} ,); 
              return card; 
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