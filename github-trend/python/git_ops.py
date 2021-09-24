from git import Repo

def download_repo(git_url: str, temp_dir: str)-> str:
    print(f'Start downloading git repo {git_url}')
    git_dir = get_git_dir(git_url)
    dir = f'{temp_dir}/{git_dir}'
    Repo.clone_from(git_url, dir)
    return dir

def get_git_dir(gitURL: str) -> str:
	segs = gitURL.split("/")
	last:str = segs[-1]
	if last.endswith(".git"): 
		return last.split(".")[0]
	return last