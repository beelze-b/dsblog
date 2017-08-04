
# coding: utf-8

# In[1]:

import pandas as pd
import numpy as np
import requests
import grequests
import os.path


# In[2]:

# auxiliary function for below
def split(theList, n):
    for i in range(0, len(theList), n):
        yield theList[i:i + n]


# In[3]:

host = "https://api.stratz.com/api/v1/match"


# In[4]:

params = {"matchId":[3304258209,3304204784], 
          "include": ["Player"],
         "gameMode": [2, 22],
         "lobbyType": [2, 7]}


# In[5]:

x = requests.get(host, params = params)


# In[9]:

x


# In[15]:

host = "https://api.stratz.com/api/v1/match/3304258209"
rs = (grequests.get(host, session=session),)
x = grequests.map(rs)
data = x[0].json()


# In[16]:

keys = data.keys()
print(keys)


# In[17]:

testsequence = [event['item'] for event in data['players'][0]['purchaseEvents']]
testitems = [45, 42, 44]
list(filter(lambda a: a not in testitems, testsequence))


# In[18]:

keys


# In[19]:

data['players'][0].keys()


# In[20]:

print(data['players'][0]['role'])
print(data['players'][0]['lane'])
print(data['players'][0]['hero'])
print(data['players'][0]['purchaseEvents'])
for player in data['players']:
    print(player['slot'])


# **Need to get steamIDs for players by going through the matchIDs of high level games**

# In[21]:

def nextBatchOfMatches(idList):
    queryParams = {"matchId":idList, 
          "include": ["Player"],
         "gameMode": [2, 22],
         "lobbyType": [2, 7]}
    x = requests.get(host, params = queryParams)
    x = x.json()
    for entry in x['results']:
        for player in entry['players']:
            playerSet.add(player['steamId'])


# In[22]:

if os.path.exists('playerset.txt'):
    playerSet = np.loadtxt('playerset.txt', dtype=np.int64)
else:
    playerSet = set()
    allMatches = np.loadtxt('matchids.txt', dtype=np.int64)
    allMatches = split(allMatches, 10)
    for matches in allMatches:
        nextBatchOfMatches(matches)
    playerSet = list(playerSet)
    np.savetxt('playerset.txt', playerSet, delimiter=',', fmt = "%.0f")


# With player set, we can obtain matches. For the players above, get the matches these players play and also categorize the players and heroes

# In[23]:

# 8 is jug and 74 is invoker and 106 is ember
def playerMatchObtain(player, heroList = [8, 74, 106]):
    host = "https://api.stratz.com/api/v1/match"
    queryParams = {"steamId": player, 
                   "heroId": heroList,
                   "gameMode": [2, 22],
                   "lobbyType": [2, 7],
                   "gameVersion": 79, 
                  "take": 100}
    x = requests.get(host, params = queryParams)
    x = x.json()
    for entry in x['results']:
        matchSet.add(entry['id'])


# In[24]:

matchSet = set()
if os.path.exists('matchset.txt'):
    matchSet = np.loadtxt('matchset.txt', dtype=np.int64)
else:
    allplayers = np.loadtxt('playerset.txt', dtype=np.int64)
    for player in allplayers:
        playerMatchObtain(player, heroList=[74])
    matchSet = list(matchSet)
    np.savetxt('matchset.txt', matchSet, delimiter=',', fmt = "%.0f")


# In[25]:

len(playerSet)


# In[26]:

len(matchSet)


# ** Time to Parse these Matches**
# Need to make sure to remove tps and wards from buy list as you can buy them at any time depending on the situation

# In[27]:

prototypeDataframe = {
    "invokerItems": [], "invokerSideRadiant": [], 
    "invokerSideHero1": [], "invokerSideHero2": [], 
    "invokerSideHero3": [], "invokerSideHero4": [],
    "otherSideHero1": [], "otherSideHero2": [],
    "otherSideHero3": [], "otherSideHero4": [],
    "otherSideHero5": [], 
    "invokerSideHero1Steam": [], "invokerSideHero2Steam": [], 
    "invokerSideHero3Steam": [], "invokerSideHero4Steam": [],
    "otherSideHero1Steam": [], "otherSideHero2Steam": [],
    "otherSideHero3Steam": [], "otherSideHero4Steam": [],
    "otherSideHero5Steam": [], "matchId": []
}
session = requests.Session()
def analyzeMatchesForFeatures(matchList, heroesToTestItems = [8, 74, 106], ):
    hostInitial = "https://api.stratz.com/api/v1/match/"
    rs = (grequests.get(hostInitial + str(match), session=session) for match in matchList)
    x = grequests.map(rs, size = 50)
    for matchUnparsed in x:
        if not matchUnparsed.ok:
            matchUnparsed.close()
            continue
        match = matchUnparsed.json()
        radiantHeroes = []
        direHeroes = []
        radiantSteams = []
        direSteams = []
        invokerSideRadiant = None
        for player in match['players']:
            # we have the invoker
            if player['hero'] == 74:
                invokerSideRadiant = player['slot'] < 5
                invokerItems = [event['item'] for event in player['purchaseEvents']]
            else:
                if player['slot'] < 5:
                    radiantHeroes.append(player['hero'])
                    radiantSteams.append(player['steamId'])
                else: 
                    direHeroes.append(player['hero'])
                    direSteams.append(player['steamId'])
        invokerWon = invokerSideRadiant == match['didRadiantWin']
        if not invokerWon:
            matchUnparsed.close()
            continue
        prototypeDataframe['matchId'].append(match['id'])
        prototypeDataframe['invokerSideRadiant'].append(invokerSideRadiant)
        prototypeDataframe['invokerItems'].append(invokerItems)
        if invokerSideRadiant:
            invokerSideHeroes = radiantHeroes
            invokerSideSteams = radiantSteams
            otherSideHeroes = direHeroes
            otherSideSteams = direSteams
        else:
            invokerSideHeroes = direHeroes
            invokerSideSteams = direSteams
            otherSideHeroes = radiantHeroes
            otherSideSteams = radiantSteams
        for ix1 in range(5):
            prototypeDataframe['otherSideHero' + str(ix1+1)].append(otherSideHeroes[ix1])
            prototypeDataframe['otherSideHero' + str(ix1+1) + 'Steam'].append(otherSideSteams[ix1])
            if ix1 != 4:
                prototypeDataframe['invokerSideHero' + str(ix1+1)].append(invokerSideHeroes[ix1])
                prototypeDataframe['invokerSideHero'+ str(ix1+1) + 'Steam'].append(invokerSideSteams[ix1])
        matchUnparsed.close()
            
def GetHeroHistory(steamid, heroid):
    host = "https://api.stratz.com/api/v1/match"
    queryParams = {"steamId": steamid, 
                   "heroId": [heroid],
                   "Player": "Single",
                   "gameMode": [2, 22],
                   "lobbyType": [2, 7],
                   "gameVersion": 79, 
                  "take": 100}
    x = requests.get(host, params = params)
    x = x.json()
    features = {'kills': [], 'deaths': [], 'lastHits': [], 'gold': [], 'obs': [], 'sens': []}
    for match in x['results']:
        pass


# In[28]:

if not os.path.exists('invokerMatchFeatureSetPlus.csv'):
    analyzeMatchesForFeatures(matchSet)
    invokerMatchFeatureSetPlus = pd.DataFrame(prototypeDataframe)
    invokerMatchFeatureSetPlus.to_csv('invokerMatchFeatureSetPlus.csv', index=False)

