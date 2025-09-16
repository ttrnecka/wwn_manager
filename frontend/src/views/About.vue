<template>
  <div>
    <h2 class="blue">Goal</h2>
    <p>
      Main goal of the WWN manager is to tackle the issue of identifying the hosts in SAN environment
      with the secondary goal of being able to handle multitenant environments with potentially duplicate WWN.
    </p>

    <h2 class="blue">Design</h2>
    <p>
      To identify the hostname, the idea is to use SAN zoning, alias or WWN data. Every WWN in SAN may be associated
      with mutliple zones and preferably one alias. Zones and aliases follow naming convention that helps to identify
      the host the WWN belongs to. The naming convention can be different for different tenants or even not being followed
      for every zone/alias within the tenant. The tool uses a set of rules to identify the host the WWN belongs to. That allows 
      the tool to generate a list of WWNs and hostnames to be used in external tool to map WWNs to hosts.
    </p>

    <h3 class="blue">WWN source</h3>
      WWNs are imported to the tool in specific format containig following fields:
    <ul>
    <li>Customer</li>
    <li>WWN</li>
    <li>Zone</li>
    <li>Alias</li>
    <li>Existing hostname</li>
    <li>Flag if the WNW is loaded by CSV in the external tool or detected</li>
    <li>WWN set</li>
    </ul>

    WWN set:
    <ul>
    <li>1 - WWN loaded from SAN configuration with zone and alias data</li>
    <li>2 - WWN manually loaded into the source system, without zone and alias data</li>
    <li>3 - WWN automatically discovered by the source system, without zone and alias data</li>
    </ul>
    <p>
      Set 1 records come with zone and alias info, while set 2 and 3 do not. If set 1 WWN belongs to more zones/aliases they should be 
      included on a separate lines in the import file. The same WWN can be imported multiple times from the SAN configuration (set 1), 
      while it can be imported only once as manually loaded (set 2) or automatically discovered (set 3).
    </p>

    <h3 class="blue">Rules</h3>
    <p>There are different type of rules that control how the host name detection for WWN is handled</p>
    <ul>
    <li>range rules</li>
    <li>host rules</li>
    <li>reconciliation rules</li>
    </ul>

    <h4 class="blue">Range Rules</h4>
    <p>
      Purpose of the range rules is to separate host WWNs from WWNs used for array, backup or other non-host purposes. They are
      applied only on global level. They use regular expresion to match given WWN to one type of array, backup, host or other. The
      range rules are applied first. Range rules are tested based on the configured order with first one matching being the rule to decide
      what type the WWN is. WWNs without no range rule applied are marked with type <i>Unknown</i>.
      <br/><br/>
      Range rules types are:
      </p>
      <ul>
        <li><b>WWN Range - Array</b> - matching WWN is considered storage array</li>
        <li><b>WWN Range - Backup</b> - matching WWN is considered backup device</li>
        <li><b>WWN Range - Host</b> - matching WWN is considered host HBA</li>
        <li><b>WWN Range - Other</b> - matching WWN is considered non host HBA, use if not sure if array or backup</li>
      </ul>

    <h4 class="blue">Host Rules</h4>
    <p>
      Host rules are applied only to <b>Host</b> type WWNs and use regular expression with match group(s) to pluck the subset of the zone or alias name
      as the hostname for the given WWN. Zone and alias based rules need to contain match group(s) and specify which match 
      group denotes the host. Alternatively the whole WWN can be used to map it to a given hostname. The host rules
      are applied after range rules. Host rules can be either global or dedicated to a tenant. They are tested in the configured order with
      tenant rules being applied before global rules and first one matching being used to decode the hostname.
    <br/><br/>
      Host rules types are:
    </p>
    <ul>
      <li><b>Alias</b> - use regular expression with match group on WWN alias to decode host</li>
      <li><b>WWN</b> - user exact WWN to match to host name. Use <b>Comment</b> field as the hostname</li>
      <li><b>Zone</b> - use regular expression with match group on WWN zone to decode host</li>
    </ul>
    <h4 class="blue">Reconciliation Rules</h4>
    <p>
      These are special rules to handle duplicate WWNs across multiple tenants or handle name missmatch between
      current hostname coming from the import and hostname decoded using the host rules. These rules are applied after host rules. 
      Additionaly there is set of default reconciliation rules that are always applied. Reconciliation rules can be either global 
      or dedicated to a tenant. The order in which they are applied is - default rules, tenant rules and global rules. They are applied as long as
      the record still requires reconciliation. If the record still requires reconciliation after all rules are applied, then manual 
      reconciliation needs to be done in the interface.
    </p>

    <h5 class="blue">Default Reconciliation Rules</h5>
    <ul>
      <li>
        <b>Rule 1</b> - if WWN is automatically discovered and it belongs to mutliple customers, then the automatically disovered record 
        will be included in the standard export. All other records with the same WWN are automatically included in the override export. 
        Additionaly, if any record's decoded hostname matches the automaticaly discovered record hostname, then such record is <b>ignored</b>
        from any export to avoid duplicate data.
      </li>
      <li>
        <b>Rule 2</b> - if WWN is manually loaded and it belongs to multiple customers, if any record's decoded hostname matches its hostname, or part of it,
        then the manually loaded WWN is included in the standard export, the matching record is <b>ignored</b> from any export to avoid duplicate data
        and all other records are included in the override export.
      </li>
      <li>
        <b>Rule 3</b> - same scenario as <b>Rule 2</b> but all hostnames are unique, then manual reconciliation is required, to decide which record goes
        to which export.
      </li>
      <li>
        <b>Rule 4</b> - if WWN is loaded from SAN configuration and it belongs to mutliple customers where are all loaded from SAN, then manual 
        reconciliation is needed to decide which record goes to which export.
      </li>
    </ul>
    <h5 class="blue">Custom Reconciliation Rules</h5>
    <p>
      Custom reconciliation rules are used to resolve records that still require manual reconciliation after the default
      reconciliation rules had been applied. They are:
    </p>
    <ul>
      <li>
        <b>WWN Primary Customer</b> - maps WWN to a customer record to be used in standard export. Creating this rule for 1 record will cause all
        other records with the same WWN to go to override export. The interface provides reconciliation pop up to select the primary customer and will
        autogenerate this rule based on that selection.
      </li>
      <li>
        <b>Ignore Loaded Host</b> - regular expression that causes the loaded hostname to be ignored when checking missmatch between
        loaded and decoded host. It is used to tell to export the record with decoded hostname despite the missmatch. Counter rule to tell to export 
        the WWN with the loaded hostname is <b>host rule</b> of <b>WWN</b> type that maps the WWN to loaded hostname. The tool allows one click setup of both these rules.
        There are 2 formats for the regular expression field here:
        <ul>
          <li>
            &lt;host_regexp&gt; - this will ignore all loaded hosts matching that regexp if reconciliation is required
          </li>
          <li>
            &lt;host_regexp&gt;,&lt;WWN&gt; - this will ignore loaded hosts matching that regexp only if the record WWN equals to the WWN
          </li>
        </ul>
      </li>
    </ul>

    <h2 class="blue">Usage</h2>

    Here is the basic usage loop for the tool:

    <ul>
      <li>
        Import entries
      </li>
      <li>
        Import rules
      </li>
      <li>
        Update range rules
      </li>
      <li>
        Reconcile records
      </li>
      <li>
        Delete records
      </li>
      <li>
        Review all changes
      </li>
      <li>
        Update host rules
      </li>
      <li>
        Save records
      </li>
      <li>
        Export the Host WWNs and Override WWNs
      </li>
    </ul>

  <h3 class="blue">Import Entries</h3>
  <ul>
    <li>Locate the <b>Host WWN Identification Feed</b> report in the reporting tool</li>
    <li>Run it and export it as CSV</li>
    <li>Click <b>Import Entries</b> button on the <b>Global</b> page and select the csv file</li>
  </ul>

  <h3 class="blue">Import Rules</h3>

  Follow this if you do not have latest rules set already:

  <ul>
    <li>Download the rules.csv from the repository (location TBD)</li>
    <li>Click <b>Import Rules</b> button on the <b>Global</b> page and select the csv file</li>
  </ul>

  <h3 class="blue">Update Range Rules</h3>

  After new import it is required to verify that there are no WWNs with <b>Unknown</b> type:

  <ul>
    <li>Click <b>Apply Rules</b> button on the <b>Global</b> page</li>
    <li>Type <b>Unknown</b> in the search box on the <b>Global</b> page</li>
    <li>If entries are listed, try to identify if they are backup, array, host or others, based on the zones and aliases they are in</li>
    <li>Once done, either update existing Range Rules or create new one</li>
      <ul>
        <li>Expand the <b>Range Rules</b> section on the <b>Global page</b></li>
        <li>Update regular expresion for existing rule or click <b>Add Rule</b> button at the bottom of the table and fill input appropriately</li>
      </ul>
    <li>Repeat all steps until there is no <b>Unknown</b> record</li>
  </ul>

  <h3 class="blue">Reconcile Records</h3>

  <ul>
    <li>Click <b>Reconcle Only</b> checkbox on the <b>Global</b> page</li>
    <li>For every displayed line, do one of the following</li>
      <ul>
        <li>Click <button class="btn btn-outline-primary btn-sm"><i class="bi bi-arrow-bar-up" role="button"></i></button> next to a hostname to select that hostname for the given WWN</li>
        <li>Click <b>Reconcile</b> button to bring up modal and select the primary customer (record that will be exported)</li>
      </ul>
    <li>Repeat all steps until there is no entry in the list</li>
  </ul>

  <h3 class="blue">Delete Records</h3>

  Follow this if you want to delete any record so it is not exported again:

  <ul>
    <li>In <b>Global</b> or <b>Customer</b> page find the WWN and click <button class="btn btn-outline-danger btn-sm"><i class="bi bi-trash text-danger" role='button'></i></button> to delete the record.</li>
    <li>Once the record is deleted you can revert it in the <b>Summary</b> view in the <b>Deleted WWN Records</b> or <b>Deleted Override WWN Records</b> sections by clicking <button class="btn btn-outline-primary btn-sm">
      <i class="bi bi-chevron-up" role='button'></i></button> button.</li>
  </ul>

  <h3 class="blue">Review All Changes</h3>

  <ul>
    <li>Navigate to the <b>Summary</b> page</li>
    <li>Review all sections - especially the new, changed and deleted sections</li>
    <ul>
      <li><b>Unaltered WWN records</b> - records where host to WWN mapping did not change after the rules have been applied</li>
      <li><b>New WWN records</b> - records where WWN to host mapping currently does not exist in the source tool</li>
      <li><b>Changed WWN records</b> - records where host is different than host already in the tool after all rules have been applied</li>
      <li><b>Deleted WWN records</b> - records that exist in the source but have been marked for deletion in the next export</li>
      <li><b>Unaltered Override WWN records</b> - duplicate override records that did not change compared to the selected snapshot</li>
      <li><b>New Override WWN records</b> - new duplicate records that will go to override table, compared agaist the selected snapshot</li>
      <li><b>Changed Override WWN records</b> - changed duplicate records in override table where the host is different than in selected snapshot </li>
      <li><b>Deleted Override WWN records</b> - deleted override records, compared to the selected snapshot</li>
      <li><b>Ignored WWN records</b> - records that have both same WWN and similar decoded hostname with another WWN record, where the other record has been prioritized by the default reconciliation rules</li>
    </ul>
  </ul>

  <h3 class="blue">Update Host Rules</h3>

  If, after reviewing the changes, some of the hosts are not captured properly create or update tenant host rules.

  <ul>
    <li>Navigate to <b>Customers</b> section and select the customer you want to edit</li>
    <li>Expand the <b>Host Rules</b> section</li>
    <li>Update regular expresion for existing rule or click <b>Add Rule</b> button at the bottom of the table and fill input appropriately</li>
      <ul>
        <li><b>Regexp</b> field must contain regular expression with capture group</li>
        <li><b>Group</b> should contain number specifying the capture group that contains the hostname</li>
        <li>Make sure the more specific rules are lower order that more general rules</li>
      </ul>
    <li>Click <b>Save</b> the rules</li>
    <li>Click <b>Apply Rules</b> and review the changes</li>
  </ul>

  <h3 class="blue">Save Records</h3>
  
  After all Unknown records and records needing reconciliation have been addressed the records can be saved. This will
  create a snapshot of the records for future reference.

  <ul>
    <li>Navigate to <b>Summary</b> section and click the <b>Save</b> button</li>
    <li>Provide optional comment to easier reference the snapshot and click <b>Save</b> again!</li>
    <li>If there are Unknown or not reconciled records the save will not go through and notification is displayed</li>
  </ul>

  <h3 class="blue">Export the Host WWNs and Override WWNs</h3>
    <ul>
      <li>Navigate to <b>Summary</b> section and click the <b>Export Host WWNs</b> or <b>Export Override WWNs</b> button</li>
      <li>Select record version to export</li>
      <li>Use these files to import the WWNs into the reporting tool.</li>
      <li>If there are Unknown or not reconciled records wont be visible</li>
      <li>The host names are exported as lower case, the source tool uses case-insensitive matching when checking host existence</li>
    </ul>
  </div>
</template>

<style>
body {
  font-size:15px;
}

</style>
