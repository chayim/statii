from github import Github
import sys
from optparse import OptionParser
from colours import Colours
import yaml


class PRData:

    # ideally title, author, url, state mergeable, repo
    def __init__(self, **kwargs):
        for k, v in kwargs.items():
            setattr(self, k.upper(), v)


IssueData = PRData


class OrgRepoData:

    def __init__(self, iam, token, orgname, labelfilter=""):
        self.GIT = Github(token)
        self.IAM = iam
        ghorg = self.GIT.get_organization(orgname)
        self.REPOS = ghorg.get_repos()
        self.LABELFILTER = labelfilter

    def _collect(self):
        prs = self.my_pull_requests()
        issues = self.my_issues()
        return prs, issues

    def my_pull_requests(self):
        prs = []
        for repo in self.REPOS:
            pulls = repo.get_pulls(state='open', sort='updated-desc')
            for pr in pulls:
                author = pr.user.login

                if author != self.IAM:
                    continue

                # reformat time
                title = pr.title
                url = pr.html_url
                state = pr.state
                mergeable = pr.mergeable
                repo = repo
                prs.append(PRData(title=title, url=url, state=state,
                          mergeable=mergeable, repo=repo.full_name.split("/")[-1]))
            return prs

    def my_issues(self):
        issues = []
        for repo in self.REPOS:
            ghissues = repo.get_issues(assignee=self.IAM, state='open', sort='updated-desc')
            for i in ghissues:
                breakout = False
                for l in i.labels:
                    if l.name == self.LABELFILTER:
                        breakout = True
                if breakout:
                    continue

                pr = i.pull_request
                title = i.title
                last_update = i.updated_at.strftime("%Y-%m-%d %H:%M")
                issues.append(IssueData(title=title, last_update=last_update, pr=pr,
                              repo=repo.full_name.split("/")[-1], url=i.html_url))
        return issues

class Plugin:

    CONFIGVAR = "github"

    def __init__(self, conffile, notify=False):
        self.NOTIFY = notify

        self.CFG = yaml.load(conffile, Loader=yaml.FullLoader).get(self.CONFIGVAR)

    def run(self):
        for o in self.CFG.get("organizations"):
            try:
                org = OrgRepoData(self.CFG.get("username"), self.CFG.get("token"), o, self.CFG.get("labelfilter"))
                prs, issues = org._collect()
            except:
                continue
            for p in prs:
                type = "PR"
                msg = f"{Colours.OKGREEN}{p.REPO} {type}:{Colours.ENDC} {Colours.BLUE}{p.TITLE}{Colours.ENDC} {Colours.WHITE}{p.URL}{Colours.ENDC}"
                print(msg)
            for p in issues:
                type = "ISSUE"
                msg = f"{Colours.OKGREEN}{p.REPO} {type}:{Colours.ENDC} {Colours.BLUE}{p.TITLE}{Colours.ENDC} {Colours.WHITE}{p.URL}{Colours.ENDC} {Colours.BLUE}{p.LAST_UPDATE}{Colours.ENDC}"
                print(msg)
