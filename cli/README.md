# STAMETS CLI

A CLI interface for aggregating results produced by STAMETS.

Aim it a directory containing logs/result diagnostics containing printed STAMETS results. To aggregate PTA metrics, give the ``-pta`` flag. To aggregate call graph metrics use ``-cg``

Example:
```
stamets -dir ./foo/bar -pta -cg
```
