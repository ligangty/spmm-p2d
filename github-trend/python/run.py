from parse_html import parseContent
from get_content import get_github_trending
from git_ops import download_repo
from repeated_timer import RepeatedTimer
import time
import os
import shutil
import tempfile

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
    temp_dir = tempfile.mkdtemp(prefix='trend-')
    repoLocalPath = download_repo(gitURL, temp_dir)
    
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
    
    print(f'Cleaning temp dirs {temp_dir}')
    shutil.rmtree(temp_dir)
    
if __name__=="__main__":
    run(3) # run first time
    rt = RepeatedTimer(30, run, 3) # it auto-starts, no need of rt.start()
    try:
        time.sleep(3600) # your long-running job goes here...
    finally:
        rt.stop() # better in a try/finally block to make sure the program ends!