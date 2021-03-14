#!/usr/bin/env python

import datetime
import requests
from colours import Colours
import yaml
import jira
from colours import Colours


class Plugin:

    CONFIGVAR = "jira"

    def __init__(self, conffile, notify=False):
        self.NOTIFY = notify

        self.CFG = yaml.load(conffile, Loader=yaml.FullLoader).get(self.CONFIGVAR)

    def run(self):
        j = jira.JIRA(self.CFG.get("server"), basic_auth=(self.CFG.get("username"), self.CFG.get("token")))
        issues = j.search_issues(self.CFG.get("query"))
        for i in issues:
            url = f"{self.CFG.get('server')}/browse/{i.fields.project}-{i.id}"
            msg = f"{Colours.OKGREEN}{i.id}:{Colours.ENDC} {Colours.BLUE}{i.fields.summary}{Colours.ENDC} {Colours.WHITE}{url}{Colours.ENDC} {Colours.BLUE}{i.fields.updated}{Colours.ENDC}"
            print(msg)
