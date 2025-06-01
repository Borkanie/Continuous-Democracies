# Continuous-Democracies
This project will allow us to see what the fuck the parlament is doing on our money at all times.

## Arhitecture

[Overall diagram](https://app.diagrams.net/#G1IoOFjcHjG2b3ap5Qtk5loAzpLCS3u-wf#%7B"pageId"%3A"tO0PR2FBSI-hfwW1G_-8"%7D)

## DataBase

It will use PostGRSQL on a kubernetes cluister to make it easy to move into the cloud and have a fast reliable data structure.
The porject will use EntityFramework to run the db. This feature also comes with a caching system that will be sufficient for our usecases.
In order to deploy changes done to model project to db isntance we have to run the follwoing commands in developer console in the folder of the DBManager (ParlimentMonitor.DataBaseConnector):
```
dotnet ef migrations add SomeName
dotnet ef database update
```
This will basically create a new commit and deploy it on the database. The commit can be seen in the itnermidiate file in the 'Migrations' direcotry inside the project. There the framework converts the changes to C# statements and than we run that on our DB after connecting to it.
Connection string will be saved in "Secrets.cs" file.
