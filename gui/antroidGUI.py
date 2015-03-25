#!/usr/bin/env python2.7
# -*- coding: UTF-8 -*-

""" Gui for project antroid """

# imports
from sys import stdin
from time import sleep
from threading import Thread
from Tkinter import Tk, Frame
from Tkinter import Canvas, PanedWindow, Spinbox
from Tkinter import Label, Listbox, Scrollbar, StringVar
from Tkinter import HORIZONTAL, VERTICAL, ALL, RIGHT, TOP, LEFT

# constants
GRASS = 0
ROCK  = 2
WATER = 4
SUGAR = 1
MILL  = 3
MEAT  = 5

STATUS_PLAY = 1
STATUS_OVER = 0

STATUS_CASE_REMEMBER = 0
STATUS_OVER_VIEW = 1

COLOR = [
        ['dark green','green'],
        ['light grey','white snow'],
        ['dim gray','gray'],
        ['dark orange','yellow'],
        ['midnight blue','blue'],
        ['seinna','tomato'],
        ]
BRAIN = [
        "Died??",
        "controlled"
        ]

class ThreadGUI(Thread):
    def __init__(self):
        Thread.__init__(self)

    def run(self):
        #là ça marche
        self.window = AntroidGUI()
        self.window.mainloop()

    def addTurn(self, turnvalue, turnnumber, readstate):
        if self.GUIstate:
            self.setreadstate(readstate)
            self.etat = self.window.addTurn(turnnumber, turnvalue)

        return self.etat

    def GUIstate(self):
        etat = self.window.state()
        if 'normal' == etat:
            return True
        else:
            return False

    def setreadstate(self, readstate):
        if self.GUIstate():
            self.window.setreadstate(readstate)

    def stop(self):
        self.window.quit()


class AntroidGUI(Tk):

    def __init__(self):
        Tk.__init__(self)
        self.protocol("WM_DELETE_WINDOW", self.myquit)
        self.etat = True
        self.readstate = True
        self.title("Antroid Map")
        self.turns = {}
        self.selectTurn = -1
        self.caselength = 10
        self.CaseMarge = 20
        self.var = StringVar()
        self.initGUI()

    def initGUI(self):
        frame = Frame(self, width =630, height = 500)
        self.panel = PanedWindow(frame, orient=HORIZONTAL)

        self.LeftFrame(self.panel)
        self.RightFrame(self.panel)
        self.panel.add(self.LFrame)
        self.panel.add(self.RFrame)
        self.panel.pack()
        frame.pack()

    def LeftFrame(self, parent):
        self.LFrame=Frame(parent, width = 500, height = 500)
        self.maptroid = Canvas(self.LFrame,bg='black',width=500,height=500)
        if self.selectTurn >= 0:
            self.printMap()

        self.maptroid.pack()
        self.LFrame.pack()

    def RightFrame(self, parent):
        self.RFrame = Frame(parent, width = 130, height = 500)
        if self.selectTurn >= 0:
            self.printInfo()
        self.FrameInfo = None
        self.updateT = True
        self.RFrame.pack()

    def setreadstate(self, readstate):
        self.readstate = readstate

    def addTurn(self, turnnumber, turnvalue):
        self.turns[turnnumber] = turnvalue
        self.updateGui(turnnumber)
        return self.etat

    def updateGui(self, turn):
        self.selectTurn = turn
        self.updateMap()
        self.updateInfo()

    def updateMap(self):
        self.maptroid.delete(ALL)
        self.printMap()

    def updateInfo(self):
        if self.FrameInfo:
            if self.updateT:
                self.frameT.destroy()
                self.printTurn()
            self.frameAnts.destroy()
            self.frameAnt.destroy()
            self.printAnts()
            self.updateT = True
        else:
            self.printInfo()

    def updateSpin_turn(self):
        turn = int(self.Spin_T.get())
        self.updateT = False
        self.updateGui(turn)

    def validateTurn(self, event):
        try:
            turn = int(self.Spin_T.get())
        except ValueError:
            turn = self.selectTurn
        if turn in self.turns.keys():
            if turn != self.selectTurn:
                self.updateT = False
                self.updateGui(turn)
        else:
            turn = self.selectTurn


    def choiceAnt(self,event):
        i = self.listbox.curselection()
        id_a = self.listbox.get(i)
        self.frameAnt.destroy()
        self.printInfoAnt(int(id_a))

    def printInfo(self):
        self.FrameInfo=Frame(self.RFrame)
        self.printTurn()
        self.printAnts()
        self.FrameInfo.pack()

    def printTurn(self):
        frameS = PanedWindow(self.FrameInfo, orient=HORIZONTAL)
        turns = Label(frameS, text="Tour :")
        self.var.set(str(self.selectTurn))
        self.Spin_T = Spinbox(frameS, values=self.turns.keys(), command =
                self.updateSpin_turn ,textvariable=self.var)
        self.Spin_T.bind('<Return>', self.validateTurn)
        turns.pack()
        self.Spin_T.pack()
        frameS.add(turns)
        frameS.add(self.Spin_T)
        frameS.pack()
        self.frameT = frameS

    def printAnts(self):
        frameAnts = Frame(self.FrameInfo)
        Text_A = Label(frameAnts, text="Fourmie :")
        s1 = Scrollbar(frameAnts)
        l1 = Listbox(frameAnts)
        id_ants = self.checkAnts()
        for i in id_ants: l1.insert(i, str(i))
        s1.config(command = l1.yview)
        l1.config(yscrollcommand = s1.set)
        l1.bind('<ButtonRelease-1>',self.choiceAnt)
        self.listbox = l1
        Text_A.pack(side = TOP)
        l1.pack(side = LEFT)
        s1.pack(side = RIGHT)
        frameAnts.pack()

        self.printInfoAnt(id_ants[0])
        self.frameAnts = frameAnts


    def printInfoAnt(self, i):
        self.frameAnt = PanedWindow(self.FrameInfo, orient=VERTICAL)
        t_Ant = Label(self.frameAnt, text="Information Ant : %d" %(i))
        (t_brain,t_energie,t_acide) = self.getInfoAnt(i)
        a_b = Label(self.frameAnt, text=t_brain)
        a_e = Label(self.frameAnt, text=t_energie)
        a_a = Label(self.frameAnt, text=t_acide)
        t_Ant.pack(side = TOP)
        self.frameAnt.add(t_Ant)
        self.frameAnt.add(a_b)
        self.frameAnt.add(a_e)
        self.frameAnt.add(a_a)
        self.frameAnt.pack()

    def printMap(self):
        turn = self.turns[self.selectTurn]
        # Information on this turn
        config = turn[0]
        Yants = turn[1]
        EAnts = turn[2]
        InitMap = turn[3]
        Cases = turn[4]

        (MaxX,MaxY,N) = InitMap
        self.MinX = 0
        self.MinY = 0
        MaxX_map = MaxX * self.caselength + (self.CaseMarge *2)
        MaxY_map = MaxY * self.caselength + (self.CaseMarge *2)
        #configure canvas
        self.maptroid.config(scrollregion=(0,0,MaxX_map,MaxY_map))
        x1 = self.CaseMarge
        y1 = self.CaseMarge
        x2 = MaxX * self.caselength + self.CaseMarge
        y2 = MaxY * self.caselength +self.CaseMarge
        self.maptroid.create_rectangle(x1, y1, x2, y2, fill="white")

        # affiche case
        for case in Cases:
            self.printCase(case)
        # affiche your ants
        for ant in Yants:
        # nb A : Your Ant 'ID X Y DX DY E A B'
            (id_a,x,y,dx,dy,e,a,b) = ant
            self.printAnt((x,y,dx,dy,b))
        # affiche enemy ants 
        for ant in EAnts:
            self.printAnt(ant)

        #to move map
        self.maptroid.bind('<ButtonPress-1>',self.grab)
        self.maptroid.bind('<B1-Motion>',self.drag)
        self.maptroid.bind('<MouseWheel>',self.mapZoom)
        self.maptroid.bind("<Button-4>", self.mapZoom)
        self.maptroid.bind("<Button-5>", self.mapZoom)

    def printCase(self, case):
        (x,y,c,s) = case
        (x1, y1, x2, y2) = self.getPoint(x, y)
        color = COLOR[c][s]
        if c%2 == 0:
            self.maptroid.create_rectangle(x1, y1, x2, y2, fill=color)
        else:
            self.maptroid.create_rectangle(x1, y1, x2, y2, fill=COLOR[GRASS][s])
            self.maptroid.create_oval(x1, y1, x2, y2, fill=color)

    def printAnt(self, ant):
        (x,y,dx,dy,brain) = ant
        (x1, y1, x2, y2) = self.getPoint(x, y)
        self.maptroid.create_oval(x1, y1, x2, y2, fill="red4")

    def getPoint(self, x, y):
        x1 = (x - self.MinX) * self.caselength + self.CaseMarge
        y1 = (y - self.MinY) * self.caselength + self.CaseMarge
        x2 = x1 + self.caselength
        y2 = y1 + self.caselength
        return (x1, y1, x2, y2)

    def grab(self,event):
        self._y = event.y
        self._x = event.x

    def drag(self, event):
        self.maptroid.yview('scroll',self._y-event.y,'units')
        self.maptroid.xview('scroll',self._x-event.x,'units')
        self._y = event.y
        self._x = event.x

    def mapZoom(self, event):
        # respond to Linux or Windows wheel event
        if event.num == 5 or event.delta == -120:
            self.caselength -= 5
            self.caselength = max(self.caselength, 10)
        if event.num == 4 or event.delta == 120:
            self.caselength += 5
            self.caselength = min(self.caselength, 40)
        self.CaseMarge = self.caselength
        self.updateMap()

    def getInfoAnt(self,id_ant):
        turn = self.turns[self.selectTurn]
        YAnt = turn[1]
        B_Text = "Brain :"
        E_Text = "Energie :"
        A_Text = "Acide :"
        for ant in YAnt:
            (id_a,x,y,dx,dy,e,a,b) = ant
            if b != 1 :
                b = 0
            if(id_ant == id_a):
                B_Text = "Brain : " + BRAIN[b]
                E_Text = "Energie : %d" % (e)
                A_Text = "Acide : %d" % (a)

        return (B_Text,E_Text,A_Text)

    def checkAnts(self):
        turn = self.turns[self.selectTurn]
        YAnt = turn[1]
        ants_id = []
        for ant in YAnt:
            (id_a,x,y,dx,dy,e,a,b) = ant
            ants_id.append(id_a)
        return ants_id

    def myquit(self):
        self.etat = False
        if not self.readstate:
            self.quit()
