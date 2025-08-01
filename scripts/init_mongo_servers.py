#!/usr/bin/env python3
"""
MongoDB Replica Set Initialization Script
This script initializes MongoDB replica set configuration.
"""
import sys

import pymongo
import yaml

CONFIG_FILE = "mongo_servers.yml"


def load_config(config_file):
    """Load MongoDB server configuration from YAML file"""
    with open(config_file, "r", encoding="utf-8") as f:
        return yaml.safe_load(f)


def test_connection(server):
    """Test connection to a MongoDB server"""
    host = server["host"]
    port = server.get("port", 27017)
    
    # For now, connect without authentication
    uri = f"mongodb://{host}:{port}/admin?directConnection=true"
    
    try:
        client = pymongo.MongoClient(uri, serverSelectionTimeoutMS=5000)
        client.admin.command("ping")
        print(f"Connected to {host}:{port} successfully.")
    except (
        pymongo.errors.ConnectionFailure,
        pymongo.errors.ServerSelectionTimeoutError,
    ) as e:
        print(f"Error connecting to {host}:{port}: {e}")
    finally:
        client.close()


def init_primary(server):
    """Initialize replica set on the primary server"""
    host = server["host"]
    port = server.get("port", 27017)
    
    # For now, connect without authentication
    uri = f"mongodb://{host}:{port}/admin?directConnection=true"
    
    try:
        client = pymongo.MongoClient(uri, serverSelectionTimeoutMS=5000)

        # Configure replica set with three members
        rs_config = {
            "_id": "rs0",
            "members": [
                {"_id": 0, "host": "mongo-0:27017"},
                {"_id": 1, "host": "mongo-1:27017"},
                {"_id": 2, "host": "mongo-2:27017"},
            ],
        }
        try:
            client.admin.command("replSetInitiate", rs_config)
            print(f"Replica set initiated on {host}:{port}.")
        except pymongo.errors.OperationFailure as e:
            print(f"Replica set initiation error (may be already initiated): {e}")
    except (
        pymongo.errors.ConnectionFailure,
        pymongo.errors.ServerSelectionTimeoutError,
    ) as e:
        print(f"Error connecting to {host}:{port}: {e}")
        sys.exit(1)
    finally:
        client.close()


def main():
    """Main function to test connections and initialize replica set"""
    config = load_config(CONFIG_FILE)

    # Test connections to all servers
    for server in config["servers"]:
        test_connection(server)

    # Initialize replica set on the first server (primary)
    init_primary(config["servers"][0])


if __name__ == "__main__":
    main()
