# Rebase notes

    git remote add upstream https://github.com/golang/go.git
    git fetch upstream --tags
    git checkout -b temp-go1.xx.y <upstream SHA for 1.xx.y release>
    git subtree split --prefix=src/crypto/tls -b tls-go1.xx.y
    git checkout <your-branch>
    git rebase tls-go1.xx.y
