#!/usr/bin/env python

import datetime
import requests
from colours import Colours
import yaml


class Plugin:

    CONFIGVAR = "circleci"

    def __init__(self, conffile, notify=False):
        self.NOTIFY = notify

        self.CFG = yaml.load(conffile, Loader=yaml.FullLoader).get(self.CONFIGVAR)
        self.HEADER = {'authorization': 'Circle-Token %s' % self.CFG.get("token")}

    def _formatter(self, results):
        branches = {}
        for r in results:
            branch = r.get("branch")
            if branch in branches.keys():
                continue

            status = r.get("status")
            try:
                when = datetime.datetime.strptime(r.get("committer_date"), "%Y-%m-%dT%H:%M:%S.000Z")
            except (TypeError, ValueError):  # duplicate run
                continue

            if when < datetime.datetime.now() - datetime.timedelta(days=6):
                continue

            url = r.get("build_url")
            branches[branch] = {'url': url, 'when': when, 'status': status}
        return branches

    def run(self):
        d = self.CFG.get("pairs")
        results = {}
        base_url = "https://circleci.com/api/v1.1/project/github"
        for key, val in d.items():
            for project in val:
                url = "{}/{}/{}?limit=100".format(base_url, key, project)
                r = requests.get(url, headers=self.HEADER)
                if r.status_code != 200:
                    results[project] = None
                    continue
                else:
                    results[project] = self._formatter(r.json())
        return self._print_results(results)

    def _print_results(self, results):
        for key, val in results.items():
            for k, v in val.items():
                if v['status'] == "success":
                    continue

                if v['status'] == "failed":
                    colour = Colours.FAIL
                else:
                    colour = Colours.GREY
                if v['status'] != 'failed':
                    print(f"{colour}{key} [{k}]{Colours.ENDC}: {Colours.WHITE}{v['url']}{Colours.ENDC} {v['when']}")
                else:
                    print(f"{colour}{key} [{k}]{Colours.ENDC}: {Colours.FAIL}{v['url']}{Colours.ENDC} {v['when']}")
        print()
