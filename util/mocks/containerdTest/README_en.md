SoarTest is not in an isolated test environment and depends on many technic things, such as containers, databases, SQL, etc.
I should tear down those dependencies as much as possible.
Fewer dependencies make a more stable test.

I make an essential decision to integrate Containerd with UnitTest.
Containerd will take over docker sooner or later.
I choose Containerd to keep up with technology after evaluating it.
Rename from SoarTest to ContainerdTest.