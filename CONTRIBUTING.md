# Contributing plugins

As statii processes plugins, plugin contributions must be made to the plugins folder.  A plugin consists of a struct containing its assumed yaml configuration, and a Gather function - do handle assembling messages. Additionally, a test function, ensuring that the plugin in question, can generate an array of messages, is required. 

As plugin testing is dependent on tokens, the individual plugin test, should read its associated token from an environment variable. This makes it possible to test the plugin, at least in a local setting.