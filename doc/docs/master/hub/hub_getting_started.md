# Getting Started with Hub

Hub is a **SaaS platform** that 
gives you access to all of **Pachyderm's functionalities
without the burden of deploying and maintaining it** locally
or in a third-party cloud platform. 

This section walks you through
the steps of creating a workspace in Hub so that
you do not need to worry about the underlying infrastructure
and can get started right away.

!!! Note
    Hub is the quickest and easiest way to get started with Pachyderm.
    In a few minutes, you are ready to experiment.
    Note that your first 4 hours are free. 
## Get started in 4 simple steps
![Hub Steps](../images/hub_steps.png)
## Before you start
Log in with your GitHub or Google account to start using [hub.pachyderm.com](https://hub.pachyderm.com). 
![Hub Login](../images/hub_login.png)
## 1- Create a Workspace 
Click the **Create a 4-hr Workspace** button and fill out the form.
![Hub workspace](../images/hub_create_workspace.png)

You just provisioned a one-node cluster that you can now use for
a limited time for free!

!!! Note
      While Pachyderm maintains a few clusters that are instantly
      available, none may be available during periods of high traffic. If
      your workspace is in a *starting* state, you might have to wait a few
      minutes for it to be ready.

## 2- Install pachctl
Hub enables you to access your workspace through Pachyderm 
CLI `pachctl` and the web interface called the **Console**.
Although you can perform most simple actions directly in the Console,
`pachctl` provides full functionality. Most likely, you will use
`pachctl` for any operation beyond the most basic workflow.
We recommend that you use `pachctl` for all data operations and
the Console to view your data and graphical representation of your
pipelines.

After your workspace creation, open a terminal window and [install 'pachctl'](https://docs.pachyderm.com/latest/getting_started/local_installation/#install-pachctl).

!!! Warning
    `kubectl` commands are not supported for the workspaces deployed
    on Hub.
## 3-4 Configure your Pachyderm context and login to your workspace by using a one-time authentication token
1. To configure a Pachyderm context and log in to your workspace
(i.e. have your `pachctl` point to your new workspace), click the **Connect** link on your workspace name in the Hub UI.

      ![Pachyderm workspace running](../images/hub_cluster_running.png)

      In your terminal window, copy, paste, and run the commands 1,2,3 listed in the instructions.

      ![Pachyderm workspace connect](../images/hub_cluster_connect.png)

1. Verify that you have set the correct context:

      ```shell
      $ pachctl config get active-context
      ```
      The system should return the name of your workspace.
      ```
      jolly-caribou
      ```

1. Create a test repo:

      ```shell
      $ pachctl create repo test
      ```
      ```shell
      $ pachctl list repo
      ```
      **System response**
      ```
      NAME CREATED       SIZE (MASTER) ACCESS LEVEL
      test 5 seconds ago 0B            OWNER    
      ```

1. Check the repo in your console:

      In the Hub UI, click the **...** button next to your workspace name,
      then **Console**. 
      
      Your console opens in a new window. Click **View Project**. 
      Your should see your newly created repo **test**.

## Next Step

You have successfully deployed and configured a Pachyderm
workspace in Hub.

Next, start creating your first pipelines: [Beginner Tutorial](../getting_started/beginner_tutorial.md).
