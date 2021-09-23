from sys import displayhook
from typing import Dict, List
from bs4 import BeautifulSoup

def parseContent(pageContent: str, top: int) -> List[Dict]:
    reposInfo=[]
    soup = BeautifulSoup(pageContent, 'html.parser')
    articles = soup.find_all("article", class_='Box-row')
    for article in articles:
        repoInfo={}
        urlAnchor = article.find("h1").find("a")
        if urlAnchor:
            repoInfo["url"] = urlAnchor["href"][1:]
        descrpitonPara = article.find("p")
        if descrpitonPara:
            repoInfo["description"] = descrpitonPara.get_text().strip()
        langSpan = article.find("span",attrs={"itemprop":"programmingLanguage"})
        if langSpan:
            repoInfo["language"] = langSpan.get_text()
        hotAnchors = article.find_all("a",class_='Link--muted d-inline-block mr-3')
        if hotAnchors:
            repoInfo["star"]=hotAnchors[0].get_text().strip()
            repoInfo["fork"]=hotAnchors[1].get_text().strip()
        reposInfo.append(repoInfo)
    return reposInfo[0:top]