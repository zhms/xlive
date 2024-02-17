#!/bin/bash
cd ../db
node export.js x_game
node export.js x_game_user
node export.js x_game_admin
node export.js x_game_order
node export.js x_game_statistic