#!/usr/bin/env bash

# This will replace all instances of a string in folder names, filenames,
# and within files.  Sometimes you have to run it twice, if directory names change.
# https://stackoverflow.com/questions/8905662/how-to-change-all-occurrences-of-a-word-in-all-files-in-a-directory

# Example usage:
# replace_string apple banana org_name

echo $1
echo $2
echo $3

# # replace within files
# find . \( -type d -name "vendor" -o -name ".git" -o -name "*.DS_Store"  -o -name "README.md" -o -name "*.sh" \) -prune -o -type f -exec sed -i '' "s/${1}/${2}/g" {} \;


# # rename filenames
find files/etc/$1/development -type f -name "${1}.main.ini" -execdir mv {} "${2}.main.ini" \;
find files/etc/$1/production -type f -name "${1}.main.ini" -execdir mv {} "${2}.main.ini" \;
find files/etc/$1/staging -type f -name "${1}.main.ini" -execdir mv {} "${2}.main.ini" \;
find files/etc/logrotate.d -type f -name "${1}" -execdir mv {} "${2}" \;
find files/etc/init -type f -name "org-${1}.conf" -execdir mv {} "${3}-${2}.conf" \;
find files/etc/init -type f -name "org-${1}-cron.conf" -execdir mv {} "${3}-${2}-cron.conf" \;


# #rename Directories
find files/etc \( -type d -name "vendor" -o -name ".git" -o -name "*.DS_Store"  -o -name "README.md" -o -name "*.sh" \) -prune  -o -type d -name "${1}" -execdir mv {} "${2}" \;  
find files/var/www \( -type d -name "vendor" -o -name ".git" -o -name "*.DS_Store"  -o -name "README.md" -o -name "*.sh" \) -prune -o -type d -name "${1}" -execdir mv {} "${2}" \;  


# Git init and update
# git init
# git remote update origin git@github.com:Somesh/$2