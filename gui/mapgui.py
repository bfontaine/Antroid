#!/usr/bin/env python 2.7
# -*- coding: UTF-8 -*-
import sys
import threading
import time
import Queue
import Tkinter

# imports
from sys import stdin
from antroidGUI import ThreadGUI
# constants
STATUS_PLAY = 1
STATUS_OVER = 0

class MapGUI(object):

    def __init__(self):
        self.selectTurn = 0
        self.initTurn()
        self.initpiper()
        self.startGUI()

    def initTurn(self):
        self.nextTurn = []
        self.YAnts = []
        self.EAnts = []
        self.MapTurn = []

    def initpiper(self):
        self.queue = Queue.Queue()
        self.timeout = 2
        self.blocking = (self.timeout is not None)
        self.use_stdio = False
        self.running = True

    def startGUI(self):
        self.gui = ThreadGUI()
        self.gui.start()

    def startRead(self):
        self.isNotFinish = True
        self.runninggui = True
        piper = threading.Thread(target=self.readstdin)
        piper.start()
        while self.isNotFinish and self.running:
            self.readInput()
            if self.runninggui:
                self.runninggui = self.gui.addTurn(self.nextTurn, self.selectTurn , self.isNotFinish)
                if not self.runninggui:
                    self.gui.stop()
            self.running = self.isNotFinish and self.running

        if self.runninggui:
            self.gui.setreadstate(False)

    def readstdin(self):
        time.sleep(1)
        while self.running:
            if self.use_stdio:
                line = sys.readline()
            else:
                try:
                    line = raw_input()
                except (EOFError):
                    self.running = False
                    self.gui.setreadstate(False)
                    break
            if self.running:
                self.queue.put(line)

    def collectStdin(self):
        while self.running:
            try:
                line = self.queue.get(self.blocking, self.timeout)
            except Queue.Empty:
                continue
            else:
                if line:
                    return line
                else:
                    sys.exit(1)
    def readInput(self):
        self.initTurn()

        # first Line: 'T A P S'
        self.decodeFirstLine(self.collectStdin())

        # nb A : Your Ant 'ID X Y -DX DY E A B'
        for LAnt in range(0,self.nbAnt):
            self.decodeYourAnt(self.collectStdin())
        self.nextTurn.append(self.YAnts)

        # N -> number for ant ennemy
        N = self.collectStdin()
        nbEAnt = int(N)

        # nb N : Ant Enemy : X Y DX DY B 
        for LAnt in range(0,nbEAnt):
            self.decodeEnemyAnt(self.collectStdin())
        self.nextTurn.append(self.EAnts)

        # Map W H N
        self.decodeMapInit(self.collectStdin())
        # Case : X Y C S 
        for case in range(0, self.nbCase):
            self.decodeCase(self.collectStdin())
        self.nextTurn.append(self.MapTurn)

    def decodeFirstLine(self, line):
        # first Line: T A P S
        InfoTurn = map(int, line.split(" "))
        self.selectTurn = InfoTurn[0] #This turn
        self.nbAnt = InfoTurn[1]
        self.isNotFinish = (InfoTurn[3] != STATUS_OVER)
        self.nextTurn.append(InfoTurn)

    def decodeYourAnt(self, line):
        # nb A : Your Ant 'ID X Y DX DY E A B'
        self.YAnts.append(map(int, line.split(" ")))

    def decodeEnemyAnt(self, line):
        # nb N : Ant Enemy : X Y DX DY B 
        self.EAnts.append(map(int, line.split(" ")))

    def decodeMapInit(self, line):
        # Map W H N
        mapInit = map(int, line.split(" "))
        self.nbCase = mapInit[2]
        self.nextTurn.append(mapInit)

    def decodeCase(self, line):
        # Case : X Y C S 
        self.MapTurn.append(map(int, line.split(" ")))

def main():
    mapTroid = MapGUI()
    mapTroid.startRead()

if __name__ == '__main__':
    main()
