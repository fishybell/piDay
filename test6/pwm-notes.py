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
pwm.start(50)                          # Start PWM with 0% duty cycle

# let's play a song
notes = ["C3", "C4", "C5", "C6", "C7", "C8"]
lengths = [1, .5, 1, .5, 1, .5]
for i in range(len(notes)):
    n = notes[i]
    l = lengths[i]
    print(n)
    pwm.ChangeFrequency(int(float(notesDict[n])))
    time.sleep(l)

pwm.stop()
GPIO.cleanup()                         # resets GPIO ports used in this program back to input mode
