# Session 01
## Checklist
- [x] Adding Version Control
- [x] Try to develop a high-level understanding of ITU-MiniTwit
- [x] Migrate ITU-MiniTwit to run on a modern computer running Linux
- [ ] Share your Work on Github

## Notes
Commands from https://github.com/itu-devops/lecture_notes/blob/master/sessions/session_01/README_TASKS.md. Other commands listed below.

### Python 3.9
Upgraded Python via https://linuxhint.com/install-python-3-9-linux-mint/.

Used `alias python=python3.9`.

### Python virtual environment
Create a venv using `python -m venv .venv`

Start the venv `. .venv/bin/activate`.

Had to install pip again `curl -sS https://bootstrap.pypa.io/get-pip.py | python3.9`.

### Saving pip dependencies
(With venv active)
```
python -m pip freeze > requirements.txt
```

### Python 2 -> 3
```
python -m pip install 2to3

2to3 -w src
```
