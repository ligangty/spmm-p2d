import unittest
from parse_html import parseContent

class ParseHtmlTest(unittest.TestCase):
    def test_parse_html(self):
        content = self.__get_github_trending_local()
        reposInfo = parseContent(content, 3)
        self.assertIs(len(reposInfo),3)
        repoInfo = reposInfo[0]
        self.assertIs(len(repoInfo), 5)
        self.assertEqual(repoInfo["url"], "public-apis/public-apis")
        self.assertEqual(repoInfo["description"], "A collective list of free APIs")
        self.assertEqual(repoInfo["language"], "Python")
        self.assertEqual(repoInfo["star"], "158,252")
        self.assertEqual(repoInfo["fork"], "17,923")
        
    def __get_github_trending_local(self) -> str:
        with open("./test.html") as f:
            return f.read()