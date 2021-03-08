
PGantt
======

PGantt is a Gantt diagramming tool that integrates with Phabricator. It makes it
much easier to reason about the time constraints, milestones, and planning than
the tools available in Phabricator itself. It uses the Conduit API to access the
Phabricator data objects and stores all the planning data together with the
tickets in such a way that it is visible and editable in the Phabricator UI.

![Screenshot](https://github.com/daedaleanai/pgantt/raw/master/screenshot.png)

Building
--------

You will need Go and NodeJS installed on your machine. Once you do, simply type:

    ]==> go generate ./...
    ]==> cd cmd/pgantt 
    ]==> go build

If you're a Daedalean employee, do the following instead:

    ]==> DDLN_GANTT=1 go generate ./...
    ]==> cd cmd/pgantt 
    ]==> go build

The above will create a self-contain executable that you can move to whatever
location in your system you wish. It will not store or read any data locally
except for the configuration file described below.

Configuration and Running
-------------------------

PGantt reads the Arcanist credential data and uses it to authenticate with
Phabricator, so chances are good that you don't have to configure anything and
can start using it right off the bat. If, however, you've never used Archanist,
create the `~/.arcrc` file along the following lines:

```json
{
  "hosts": {
    "https://phabricator.yourdomain.com/api/": {
      "token": "cli-xxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    }
  },
  "pgantt": {
    "projects": ["Test", "Test 2"],
    "port": 9999
  }
}
```

You can get the API Token by visiting:
`https://phabricator.yourdomain.com/conduit/login/`.

Since the Phabricator's Conduit API does not allow for easy qurying of the
information necessary to run PGantt, PGantt fetches all the information abount
all the projects you follow at startup, and then does smart updates, polling,
and data caching. This means that the regular opration is fast, but the startup
is slow. Therefore, you may want to limit the number of projects PGantt follows
by specifying them in the config file. If you're impatient, you may also follow
in more detail what is happening by adding `-log-level Debug` to the program's
commandline.

You can simply run the PGannt executable in the terminal window:

    ==> ./pgantt

It will start serving the user interface at `http://localhost:9999` of whatever
other port you configured.

Configuring Phabricator (for administrators)
--------------------------------------------

PGantt needs the additional planning data to be stored along the tickets.
Therefore, you need to configure your Phabricator instance to keep track of
these extra data fields. To that end, go to `
Config -> Application Settings -> Mainphest -> maniphest.custom-field-definitions`
and add the following JSON definitions in the editor window and press the
`Save Config Entry` button.

```json
{
  "daedalean.scheduled": {
    "name": "Scheduled",
    "type": "bool",
    "default": false,
    "view": true,
    "edit": true
  },
  "daedalean.start_date": {
    "name": "Start Date",
    "type": "date",
    "view": true,
    "edit": true
  },
  "daedalean.duration": {
    "name": "Duration",
    "type": "int",
    "default": 0,
    "view": true,
    "edit": true
  },
  "daedalean.type": {
    "name": "Type",
    "type": "select",
    "default": "daedalean:task",
    "options": {
      "daedalean:task": "Task",
      "daedalean:milestone": "Milestone",
      "daedalean:project": "Project"
    },
    "view": true,
    "edit": true
  },
  "daedalean.progress": {
    "name": "Progress",
    "type": "int",
    "default": 0,
    "view": true,
    "edit": true
  },
  "daedalean.successors": {
    "name": "Successors",
    "type": "text",
    "default": "[]",
    "view": true,
    "edit": true
  }
}
```

Phabricator then seems to need to have these fields enabled in
`Config -> Application Settings -> Mainphest -> maniphest.fields`. Simply go there
and press the `Save Config Entry` button.

You're all set!
