#!/usr/bin/env python3
"""
MongoDB Replica Set Status Checker
This script checks the status of MongoDB replica set nodes.
"""

import sys

from pymongo import MongoClient
from pymongo.errors import PyMongoError


def get_state_name(state):
    """Convert numeric state to human-readable state name"""
    states = {
        0: "STARTUP",
        1: "PRIMARY",
        2: "SECONDARY",
        3: "RECOVERING",
        5: "STARTUP2",
        6: "UNKNOWN",
        7: "ARBITER",
        8: "DOWN",
        9: "ROLLBACK",
        10: "REMOVED",
    }
    return states.get(state, f"UNKNOWN({state})")


def check_replicaset_status():
    """Check and display replica set status"""
    try:
        # Connect to MongoDB using the primary node port with authentication
        connection_string = (
            "mongodb://mongo-0:mongo-0@127.0.0.1:27034/?replicaSet=rs0&authSource=admin"
        )
        client = MongoClient(connection_string)

        # Get replica set status from MongoDB
        status = client.admin.command("replSetGetStatus")

        # Print overall replica set status
        print(f"\nReplicaSet Name: {status.get('set', 'N/A')}")

        # Find and print status of each member
        members = status.get("members", [])
        print("\nMember Status:")
        print("-" * 50)

        for member in members:
            state = get_state_name(member.get("state"))
            print(f"Host: {member.get('name')}")
            print(f"State: {state}")
            print(f"Health: {'UP' if member.get('health') == 1 else 'DOWN'}")
            print("-" * 50)

    except PyMongoError as e:
        print(f"Error connecting to MongoDB: {e}")
        sys.exit(1)
    except (ValueError, TypeError) as e:
        print(f"An error occurred: {e}")
        sys.exit(1)
    finally:
        client.close()


if __name__ == "__main__":
    check_replicaset_status()
