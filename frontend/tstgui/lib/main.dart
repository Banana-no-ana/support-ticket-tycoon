import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;

void main() {
  runApp(const MyApp());
}


class Case {
  final int CaseID;
  final String State;     //Case State
  final String Assignee;  //Worker UID

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

Future<Case> createCase() async {
  final response = await http
      .get(Uri.parse('http://localhost:8001/case/create'), 
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
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class NewCaseCard extends StatelessWidget{

  final Future<Case> futureCase; 

  NewCaseCard({Key? key, 
      required this.futureCase}) : super(key:key);   

  @override
  Widget build(BuildContext context) {
    return Card(
        child: FutureBuilder<Case>(
        future: futureCase,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            if (snapshot.hasData) {
              return Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[Text(snapshot.data!.CaseID.toString())],
              );
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
      )); 
  }
}


class MyHomePage extends StatefulWidget {
  const MyHomePage({Key? key, required this.title}) : super(key: key);
  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  var newCases = <NewCaseCard>[]; 

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(children: [
          Column(
            children: newCases,
          ), 
          ElevatedButton(
            child: const Text('Generate Case'),
            onPressed: () {
              setState(() {
                newCases.add(NewCaseCard(futureCase: createCase()));
              });
            },
          ),
          ], 
        )
     )
    );
  }
}
