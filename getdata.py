
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


# In[5]:

params = {"matchId":[3304258209,3304204784], 
          "include": ["Player"],
         "gameMode": [2, 22],
         "lobbyType": [2, 7]}


# In[6]:

x = requests.get(host, params = params)


# In[7]:

x


# In[8]:

host = "https://api.stratz.com/api/v1/match/3304258209"
rs = (grequests.get(host),)
x = grequests.map(rs)
data = x[0].json()


# In[8]:

keys = data.keys()
print(keys)


# In[9]:

testsequence = [event['item'] for event in data['players'][0]['purchaseEvents']]
testitems = [45, 42, 44]
list(filter(lambda a: a not in testitems, testsequence))


# In[10]:

keys


# In[11]:

data['players'][0].keys()


# In[12]:

print(data['players'][0]['role'])
print(data['players'][0]['lane'])
print(data['players'][0]['hero'])
print(data['players'][0]['purchaseEvents'])
for player in data['players']:
    print(player['slot'])


# **Need to get steamIDs for players by going through the matchIDs of high level games**

# In[13]:

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


# In[14]:

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

# In[15]:

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


# In[16]:

print('Loading matches')
matchSet = set()
if os.path.exists('matchset.txt'):
    matchSet = np.loadtxt('matchset.txt', dtype=np.int64)
else:
    allplayers = np.loadtxt('playerset.txt', dtype=np.int64)
    for player in allplayers:
        playerMatchObtain(player, heroList=[74])
    matchSet = list(matchSet)
    np.savetxt('matchset.txt', matchSet, delimiter=',', fmt = "%.0f")


# In[17]:

len(playerSet)


# In[18]:

len(matchSet)


# ** Time to Parse these Matches**
# Need to make sure to remove tps and wards from buy list as you can buy them at any time depending on the situation

# In[35]:

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


# In[36]:

from scipy import stats
def GetHeroHistory(steamid, heroid):
    host = "https://api.stratz.com/api/v1/match"
    queryParams = {"steamId": steamid, 
                   "heroId": heroid,
                   "include": "Player",
                   "playerType": "Single",
                   "gameMode": [2, 22],
                   "lobbyType": [2, 7],
                   "gameVersion": 79, 
                  "take": 100}
    x = requests.get(host, params = queryParams)
    x = x.json()
    features = {'kills': [], 'deaths': [], 'assists': [], 
                'lastHits': [], 'denies': [], 'gold': [], 'xp': [],
                'heroDamage': [], 'heroHealing': []}
    # operating on the game level
    for game in x['results']:
        # duration is in seconds
        duration = game['duration']
        features['kills'].append(game['players'][0]['numKills'])
        features['deaths'].append(game['players'][0]['numDeaths'])
        features['assists'].append(game['players'][0]['numAssists'])
        features['lastHits'].append(game['players'][0]['numLastHits'])
        features['denies'].append(game['players'][0]['numDenies'])
        features['gold'].append(game['players'][0]['goldPerMinute'] * 1.0 * duration/60)
        features['xp'].append(game['players'][0]['expPerMinute'] * 1.0 * duration/60)
        features['heroDamage'].append(game['players'][0]['heroDamage'])
        features['heroHealing'].append(game['players'][0]['heroHealing'])
    return {'kills': np.mean(features['kills']), 'deaths': np.mean(features['deaths']), 
            'assists': np.mean(features['assists']), 
            'lastHits': np.mean(features['lastHits']), 'denies': np.mean(features['denies']), 
            'gold': np.mean(features['gold']), 'xp': np.mean(features['xp']),
            'heroDamage': np.mean(features['heroDamage']), 'heroHealing': np.mean(features['heroHealing'])} 


# In[28]:

GetHeroHistory(54325937, 74)


# In[6]:

print('Loading invoker data')
if not os.path.exists('invokerMatchFeatureSetPlus.csv'):
    analyzeMatchesForFeatures(matchSet)
    invokerMatchFeatureSetPlus = pd.DataFrame(prototypeDataframe)
    invokerMatchFeatureSetPlus.to_csv('invokerMatchFeatureSetPlus.csv', index=False)
else:
    invokerMatchFeatureSetPlus = pd.read_csv('invokerMatchFeatureSetPlus.csv')


# **Time to augment the invoker matches**

# In[37]:

def invokerRowAugment(row):
    data = {}
    data['matchId'] = int(row['matchId'])
    # go over the four players on invoker side and five on other side
    for ix in range(5):
        heroInfo = GetHeroHistory(row['otherSideHero' + str(ix+1) + 'Steam'], row['otherSideHero' + str(ix+1)])
        data['OSHero' + str(ix+1) + 'hero'] = int(row['otherSideHero' + str(ix+1)])
        data['OSHero' + str(ix+1) + 'kills'] = heroInfo['kills']
        data['OSHero' + str(ix+1) + 'deaths'] = heroInfo['deaths']
        data['OSHero' + str(ix+1) + 'assists'] = heroInfo['assists']
        data['OSHero' + str(ix+1) + 'lastHits'] = heroInfo['lastHits']
        data['OSHero' + str(ix+1) + 'denies'] = heroInfo['denies']
        data['OSHero' + str(ix+1) + 'gold'] = heroInfo['gold']
        data['OSHero' + str(ix+1) + 'xp'] = heroInfo['xp']
        data['OSHero' + str(ix+1) + 'heroDamage'] = heroInfo['heroDamage']
        data['OSHero' + str(ix+1) + 'heroHealing'] = heroInfo['heroHealing']
        if ix != 4:
            heroInfo = GetHeroHistory(row['invokerSideHero' + str(ix+1) + 'Steam'], 
                                      row['invokerSideHero' + str(ix+1)])
            data['ISHero' + str(ix+1) + 'hero'] = int(row['invokerSideHero' + str(ix+1)])
            data['ISHero' + str(ix+1) + 'kills'] = heroInfo['kills']
            data['ISHero' + str(ix+1) + 'deaths'] = heroInfo['deaths']
            data['ISHero' + str(ix+1) + 'assists'] = heroInfo['assists']
            data['ISHero' + str(ix+1) + 'lastHits'] = heroInfo['lastHits']
            data['ISHero' + str(ix+1) + 'denies'] = heroInfo['denies']
            data['ISHero' + str(ix+1) + 'gold'] = heroInfo['gold']
            data['ISHero' + str(ix+1) + 'xp'] = heroInfo['xp']
            data['ISHero' + str(ix+1) + 'heroDamage'] = heroInfo['heroDamage']
            data['ISHero' + str(ix+1) + 'heroHealing'] = heroInfo['heroHealing']
    return pd.Series(data)


# In[38]:

print('Loading supplemental data')
if not os.path.exists('supplementalInvokerMatchData.csv'):
    supplementalInvokerMatchData = invokerMatchFeatureSetPlus.apply(invokerRowAugment, axis=1)
    supplementalInvokerMatchData.to_csv('supplementalInvokerMatchData.csv', index=False)
else:
    supplementalInvokerMatchData = pd.read_csv('supplementalInvokerMatchData.csv')

