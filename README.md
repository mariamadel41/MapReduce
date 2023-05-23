# MapReduce Project

This project implements a MapReduce architecture for character counting in DNA sequences using Go language. It consists of a master node, which also acts as the client, and two slave nodes for parallel processing.

## Project Structure

The project repository is organized as follows:

mapreduce/
├── master/
├── slave1/
├── slave2/
└── README.md



- `master/`: Contains the code for the master node, responsible for task distribution and result aggregation.
- `slave1/`: Contains the code for the first slave node.
- `slave2/`: Contains the code for the second slave node.
- `README.md`: This file, providing an overview of the project and setup instructions.

## System Requirements

- Go language (version X.X.X)
- Operating System: Linux/MacOS/Windows

## Setup and Execution

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/mapreduce.git
