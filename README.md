# big-dumb-launcher
simple game launcher experiment in go

Goal is to have a “port” of the command line or desktop that is usable by a controller. You are encouraged to do setup via mouse and keyboard, FTP, SSH, config files, etc, but then do your actual game launching/controller/wifi/sound/power, etc with your controller. There will ideally be little to zero config files other than controller mappings and “last opened”.

**Basically: easy for users, easy for (linux) masters, no in between.**

* Implemented
    * Launch macOS/linux programs from $PATH based on your file structure and naming.
        * i.e. ares.system_name/Game Boy/Pokemon - Blue.gb => ares —system_name “Game Boy” 
        * Able to chain multiple arguments, allowing additional content to be loaded
            * i.e. woof.iwad/doom.file/SIGIL (v1.21).zip => woof -iwad doom -file “SIGIL (v1.21).zip”
* TODO
    * Get auto/ dirs working with smooth doom or a similar global doom mod
    * Rework the controls
        * Left/Right to nav folders in addition to A/B (enter/backspace)
        * Add START (or spacebar) as a command, don’t have it do anything yet (just print something).
            * Plan is to have it launch at any level of the path hierarchy.
                * ares.system_name/ => ares
                * woof.iwad/doom.file => woof -iwad doom
    * Ideally, if the first half of a folder name (“Pokemon Blue.alt/”) matches a filename  in the same directory I want it to:
        * Implement an extension ignore list before anything (mostly savefiles like .rom) 
        * Show as the file (Pokemon Blue.zip) to the user.
        * If you press start, it simply starts the zip as if nothing is unusual.
        * If the user presses right, it opens the matching folder
        * If it is a child/chained directory and NO matching file is found, start it anyway but cut out the second half
            * i.e. woof.iwad/doom.files => woof -iwad doom 
    * Begin filtering the (visual) names in the frontend and implementing the 3-column/pane layout
        * Past 3 levels, the fourth will open a smaller “window” or modal with its contents
    * Controller mapping implementation
        * Check if an xboxdrv mapping exists for the given controller’s GUID upon and prompt the user to create one if not.
            * The mapper is planned to be a separate program entirely so that other devices can use it, retro handhelds come to mind.
    * Find some way to abstract resolution/refresh/max res handling
        * gamescope is an option but i'd rather avoid it if i can (very x86 and AMD focused). This may be the one thing i'd need some kind of config file mapping for.