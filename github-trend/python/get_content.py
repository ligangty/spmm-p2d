
import requests

GITHUB_TREANDING = "https://github.com/trending"

def get_github_trending() -> str:
    print(f'Start requesting {GITHUB_TREANDING}')
    r = requests.get(GITHUB_TREANDING)
    if r.status_code==200:
        return r.text
    