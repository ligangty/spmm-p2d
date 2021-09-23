from parse_html import parseContent
from get_content import get_github_trending
from git_ops import download_repo
import os
import shutil

def run(top):
    content = get_github_trending()
    reposInfo = parseContent(content, top)
    for repoInfo in reposInfo:
        print(f'Repo:        {repoInfo["url"]}')
        print(f'Description: {repoInfo["description"]}')
        print(f'Language:    {repoInfo["language"]}')
        print(f'Star:        {repoInfo["star"]}')
        print(f'Fork:        {repoInfo["fork"]}' )
        print()
    gitURL = f'https://github.com/{reposInfo[0]["url"]}'
    repoLocalPath = download_repo(gitURL)
    
    readme = ""
    for (dir,_,names) in os.walk(repoLocalPath):
        done = False
        for path in names:
            fileName = path.split( "/")[-1]
            if fileName.lower().startswith("readme"):
                readme = os.path.join(dir, path)
                done = True
                break
        if done:
            break    
    
    with open(readme) as f:
        print(f.read())        
    
    shutil.rmtree(repoLocalPath)
    
if __name__=="__main__":
    run(3)