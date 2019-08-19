import csv

# import our frequency chart
notesList = list(csv.reader(open('notes.txt', 'rb'), delimiter='\t'))
notesDict = {}
for i in notesList:
    notesDict[i[0]] = i[1]

# setup pwm
import RPi.GPIO as GPIO   # Import the GPIO library.
import time								# Import time library
GPIO.setmode(GPIO.BOARD)  # Set Pi to use pin number when referencing GPIO pins.
GPIO.setup(12, GPIO.OUT)  # Set GPIO pin 12 to output mode.

pwm = GPIO.PWM(12, 10)   # Initialize PWM on pwmPin 100Hz frequency
pwm.start(0)                          # Start PWM with 0% duty cycle

# let's play a song
notes = ["F5", "F5", "F5", "F5", "C5", "D5", "F5", "D5", "F5"]
lengths = [.17, .17, .17, 0.93, 0.92, 0.92, 0.34, 0.17, 3]
gaps = [0.16, 0.16, 0.16, 0.08, 0.08, 0.08, 0.33, 0.16, 0]
tempo = 137.0
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
