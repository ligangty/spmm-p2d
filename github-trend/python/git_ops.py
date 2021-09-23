from git import Repo
import tempfile
import os

def download_repo(git_url: str)-> str:
    git_dir = get_git_dir(git_url)
    temp_dir = tempfile.mkdtemp(prefix='')
    dir = f'{temp_dir}/{git_dir}'
    Repo.clone_from(git_url, dir)
    return dir

def get_git_dir(gitURL: str) -> str:
	segs = gitURL.split("/")
	last:str = segs[-1]
	if last.endswith(".git"): 
		return last.split(".")[0]
	return last