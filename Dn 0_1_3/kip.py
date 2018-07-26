import re
import os
import sys
import shutil

def replace_all(text, dic):
    for i, j in dic.items():
        text = re.sub(i,j,text)
    return text

help_text = """
Usage:
kip install path/to/package
kip upgrade path/to/package
kip uninstall package_name
"""

try:
    if sys.argv[1] == "install":
        shutil.copytree(os.path.abspath(sys.argv[2]), os.path.abspath("src/"+os.path.basename(sys.argv[2])))

        with open(os.path.abspath("src/config.txt"), "a") as f:
            f.write(os.path.basename(sys.argv[2])+"\n")

        with open(os.path.abspath("src/"+os.path.basename(sys.argv[2])) + "/dependencies.txt", "r") as f:
            for dep in f.readlines():
                dep = dep.replace("\n", "")
                os.system("go get \"{}\"".format(dep))
                print("got {}".format(dep))

    elif sys.argv[1] == "uninstall":
        with open(os.path.abspath("src/config.txt"), "r") as f:
            contents = f.read()
            
        with open(os.path.abspath("src/config.txt"), "w") as f:
            if sys.argv[2] in contents:
                f.write(replace_all(contents,
                                    {r"\b{}\b\n".format(sys.argv[2]): ""}))
                shutil.rmtree(os.path.abspath("src/"+os.path.basename(sys.argv[2])))
            else:
                f.write(contents)
                raise Exception("package not installed.")
                
            f.close()
            
    elif sys.argv[1] == "upgrade":
        with open(os.path.abspath("src/config.txt"), "r") as f:
            contents = f.read()
            
        with open(os.path.abspath("src/config.txt"), "w") as f:
            if os.path.basename(sys.argv[2]) in contents:
                f.write(replace_all(contents,
                                    {r"\b{}\b\n".format(os.path.basename(sys.argv[2])): ""}))
                shutil.rmtree(os.path.abspath("src/"+os.path.basename(sys.argv[2])))
            else:
                f.write(contents)
                raise Exception("package not installed.")

            f.close()
        
        shutil.copytree(os.path.abspath(sys.argv[2]), os.path.abspath("src/"+os.path.basename(sys.argv[2])))

        with open(os.path.abspath("src/config.txt"), "a") as f:
            f.write(os.path.basename(sys.argv[2])+"\n")

        with open(os.path.abspath("src/"+os.path.basename(sys.argv[2])) + "/dependencies.txt", "r") as f:
            for dep in f.readlines():
                dep = dep.replace("\n", "")
                os.system("go get \"{}\"".format(dep))
                print("got {}".format(dep))
                

except Exception as e:
    print("\nError: " + str(e))
    print(help_text)