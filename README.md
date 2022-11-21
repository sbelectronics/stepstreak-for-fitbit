Stepstreak Calculator for Fitbit
Scott Baker, http://www.smbaker.com/

My Fitbit reset its stepstreak counter, even though I met my goal. This is known well-documented bug with many complaints on the fitbit forum.
So I wrote my own tool to compute the step streak.

Note: I'm a software engineer, a programmer with approximately 35 years of experience. This tool was written for my own personal use; I've made
no attempt to make it user friendly or configurable.

# Downloading Your Fitbit Monthly CSV data

If you go to the Fitbit website, and click through the gear icon, you'll eventually find a page called "Data Export". You can export your data, one
month at a time, and download it as a CSV file. This is invconvenient, but doable, and it's in a nice summary file.

There's also an option to "Export your Account Archive". This will take some time (mine was about an hour), but it will export years worth of data.
Unfortunately, it outputs the data in an inconvenient format, raw JSON with many different datapoints per day. It's not as easily consumed as the CSV
archive, but it is usable. I recommend downloading the individual CSV files a month at a time.

# Converting the "Account Archive" to a CSV file

I wrote a tool, in cmd/fitbit-import that will convert the ugly dump file from fitbit into a CSV file that includes the step count. I made no
attempt to fill in the rest of the data.

I did notice that some of the step counts returned by my importer were incorrect. Approximately two days underreported in ten months of data. I'm
not sure whether this is an error in my importer, or a problem with fitbit's exporter. Regardless, I RECOMMEND DOWNLOADING THE MONTHLY CSVs INSTEAD.

# Generating a step streak report from the CSV files

The source for this is in cmd/fitbit-stepstreak. Execute the command like this, assuming "exports" is where you put all your CSV files:

```bash
# If you're running Linux

linux/fitbit-stepstreak exports/*.csv

# If you're running Windows (don't include the wildcard; program will look for *.CSV in the directory you provide
# ... do make sure to include the trailing backslash)

windows\fitbit-stepstreak exports\
```

It will write output to the console like this:

```bash
2022/11/20 15:25:30 Stepstreak calculator for fitbit, by Scott Baker, http://www.smbaker.com/
2022/11/20 15:25:30 Step count of 9944 is below goal on 5/21/2021
You're on a 547 day step streak!
```
