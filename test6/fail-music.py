import csv

# import our frequency chart
notesList = list(csv.reader(open('notes.txt', 'rb'), delimiter='\t'))
notesDict = {}
for i in notesList:
    notesDict[i[0]] = i[1]
    if i[0] == "C#5":
        print(i[0])
        print(notesDict[i[0]])
    if i[1] == "554.37":
        print(i[0])
        print([ord(c) for c in i[0]])
        print([ord(c) for c in "C#5"])
        print(notesDict[i[0]])

# setup pwm
import RPi.GPIO as GPIO   # Import the GPIO library.
import time								# Import time library
GPIO.setmode(GPIO.BOARD)  # Set Pi to use pin number when referencing GPIO pins.
GPIO.setup(12, GPIO.OUT)  # Set GPIO pin 12 to output mode.

pwm = GPIO.PWM(12, 10)   # Initialize PWM on pwmPin 100Hz frequency
pwm.start(0)                          # Start PWM with 0% duty cycle

# let's play a song
notes = ["C5", "C#5", "D5", "B4", "F5", "F5", "F5", "E5", "D5", "C5", "E4", "E4", "C4"]
lengths = [0.105, 0.105, 0.25, 0.4, 0.4, 0.4, 0.56, 0.56, 0.56, 0.4, 0.4, 0.4, 0.4]
gaps = [0.02, 0.02, 3, 0.1, 0.6, 0.1, 0.1, 0.1, 0.1, 0.1, 0.6, 0.1, 0.1]
tempo = 200.0
timing = 1
for i in range(len(notes)):
    n = notes[i]
    l = lengths[i]
    g = gaps[i]
    print(n)
    pwm.start(50)                          # Start PWM with 0% duty cycle
    pwm.ChangeFrequency(int(float(notesDict[n])))
    time.sleep(timing * (l / tempo) * 60)
    pwm.start(0)                          # Start PWM with 0% duty cycle
    time.sleep(timing * (g / tempo) * 60)

pwm.stop()
GPIO.cleanup()                         # resets GPIO ports used in this program back to input mode
