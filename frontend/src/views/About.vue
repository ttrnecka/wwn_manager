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
        <b>Ignore Loaded Host</b> - regular expression with match group that causes the loaded hostname to be ignored when checking missmatch between
        loaded and decoded host. It is used to tell to export the record with decoded hostname despite the missmatch. Counter rule to tell to export 
        the WWN with the loaded hostname is <b>host rule</b> of <b>WWN</b> type that maps the WWN to loaded hostname. The tool allows one click setup of both these rules. 
      </li>
    </ul>

    <h2 class="blue">Usage</h2>

    Here is the basic usage loop for the tool:

    <ul>
      <li>
        Import latest entries from the source
      </li>
      <li>
        Import latest rule set from repository
      </li>
      <li>
        Identify any Unknow type records and if any create new range rules to tell if they are host, array, backup or other
      </li>
      <li>
        Click Apply rules
      </li>
      <li>
        Reconcile all records that need reconciliation
      </li>
      <li>
        In Summary view review all changes
      </li>
      <li>
        If needed update global or tenant host rules to capture correct name
      </li>
      <li>
        Once happy with all rules and everything is reconciled click Commit in the Summary view
      </li>
      <li>
        This will take snapshot of current data for future reference and allows to download Host WWN amd Override WWN files
      </li>
      <li>
        Download the Host WWN and Override WWN and import it into the external tool as needed.
      </li>
    </ul>
  </div>
</template>

<style>
body {
  font-size:15px;
}

</style>
