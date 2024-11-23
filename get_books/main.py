import requests
import sys
base_url = 'http://mirror.csclub.uwaterloo.ca/gutenberg/'
for index in sys.argv[1:]:
    with open(index, "r") as f:
        lines = f.readlines()
        for l in lines:
            id = '<a class="link" href="/ebooks/'
            if id in l and 'sort_order' not in l:
                key = l.lstrip(id).split()[0].rstrip("\"")
                url_suffix = '/'.join(list(key)[:-1]) + '/' + key + '/' + key + '.txt'
                full_url = base_url + url_suffix
                r = requests.get(full_url)
                if r.status_code == 200:
                    print(full_url)
                    open(f'books/{key}.txt', 'wb').write(r.content)
