import sys
import os
import time
import datetime

SCRIPT_DIR = os.path.dirname(__file__)
if SCRIPT_DIR not in sys.path:
	sys.path.insert(0, SCRIPT_DIR)

from log_writer import file_log

# Direct imports (files live in this directory)
from import_politicians import importPoliticians
from import_law import goFromLastforward

def run_cycle():
	year = datetime.datetime.now().year
	file_log(f"Orchestrator cycle start: {datetime.datetime.now().isoformat()}", alsoPrint=True)

	# 1) Import politicians only for the current year (IDs 1..333)
	#try:
	#	file_log(f"Importing politicians for year {year}", alsoPrint=True)
	#	importPoliticians(1, 333, year)
	#except Exception as e:
	#	file_log(f"Failed to import politicians: {e}", alsoPrint=True)

	# 2) Import new laws from last forward, attempt up to 100
	try:
		file_log("Calling import_law.goFromLastforward(100)", alsoPrint=True)
		goFromLastforward(100)
	except Exception as e:
		file_log(f"Failed to import laws: {e}", alsoPrint=True)

	file_log(f"Orchestrator cycle end: {datetime.datetime.now().isoformat()}", alsoPrint=True)

def main():
	file_log("Starting orchestrator (runs once every hour).", alsoPrint=True)
	try:
		while True:
			run_cycle()
			time.sleep(3600)
	except KeyboardInterrupt:
		file_log("Orchestrator stopped by user.", alsoPrint=True)

if __name__ == "__main__":
	main()

