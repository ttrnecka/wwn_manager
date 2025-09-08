<template>
  <div>
    <h2 class="blue">Goal</h2>
    <p>
      Main goal of the WWN manager is to tackle the issue of identifying the hosts in SAN environment
      with the secondary goal of being able to handle multitenant environments with potentially duplicate WWN.
    </p>

    <h2 class="blue">Design</h2>
    <p>
      To identify the hostname the idea is to use SAN zoning, alias or WWN data. Every WWN in SAN may be associated
      with mutliple zones and preferably one alias. Zones and aliases follow naming convention that helps to identify
      the host the WWN belongs to. The naming convention can be different for different tenants or even not being followed
      for every zone/alias within the tenant. The tool uses a set of rules to identify the host using that WWN. That allows 
      the tool to generate a list of WWNs and hostnames to be used in external tool for mapping WWNs to hosts.
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
    <li>1 - WWN loaded from SAN configuration</li>
    <li>2 - WWN manually loaded into the source system</li>
    <li>3 - WWN automatically loaded into the source system</li>
    </ul>
    <p>
      Set 1 records come with zone and alias info, while set 2 and 3 do not. Hence if set 1 belongs to more zones/aliases they should be 
      included on a separate lines in the import file. There can be only one set 2 or set 3 type record for any given WWN and multiple set 1 records
    </p>

    <h3 class="blue">Range Rules</h3>
    <p>
      Range rules are rules applied on global level. Their purpose is to separate host WWNs from WWNs used for array, backup or 
      other non-host purposes. They use regular expresion to match given WWN to one type of array, backup, host or other. The
      range rules are applied first. Range rules are tested based on the defined order with first one matching being the rule to apply.
    </p>

    <h3 class="blue">Host Rules</h3>
    <p>
      Host rules are applied only to host WWNs and user regular expression with match group(s) to pluck the subset of the zone or alias name
      as the hostname for the given WWN. Zone and alias based rules need to contain match group(s) and specify which match 
      group denotes the host. Alternatively the whole WWN can be used to map it to a given hostname. The host rules
      are applied after range rule. Host rules can be either global or dedicated to a tenant. They are tested in defined order with
      tenant rules being applied before global rules and first one matching being the one to apply.
    </p>
    
    <h3 class="blue">Reconciliation Rules</h3>
    <p>
      These are special rules to handle duplicate WWNs accross multiple tenants or handle name missmatch between
      current hostname and hostname decoded using the host rules. These rules are applied after host rules. Additionaly there is 
      set of default reconciliation rules that are always applied. Reconciliation rules can be either global or dedicated to a tenant.
      The order in which they are applied is - default rules, tenant rules and global rules. They are applied as long as
      the record still requires reconciliation. If the record still requires reconciliation after all rules are applied, then manual 
      reconciliation needs to be done in the interface.
    </p>

    <h4 class="blue">Default Reconciliation Rules</h4>
    <ul>
      <li>
        Rule 1 - if WWN is automatically discovered (set 3) and it belongs to mutliple customers, then it will be included in the standard export. 
        All other records with the same WWN are automatically included in the override export. Additionaly, if any record's decoded 
        hostname matches the automaticaly discovered record hostname, then such record is ignored from any export to avoid duplicate data.
      </li>
      <li>
        Rule 2 - if WWN is manually imported (set 2) and it belongs to multiple customers, if any record's decoded hostname matches its hostname,
        then the manually imported WWN is included in the standard export, the matching record not exported and all other records are 
        included in the override export.
      </li>
      <li>
        Rule 3 - same scenario as Rule 2 but all hostnames are unique, then manual reconciliation is required, to decide which record goes
        to which export.
      </li>
      <li>
        Rule 4 - if WWN is loaded from SAN (set 1) and it belongs to mutliple customers where are all loaded from SAN, then manual 
        reconciliation is needed to decide which record goes to which export.
      </li>
    </ul>
    <h4 class="blue">Custom Reconciliation Rules</h4>
    <ul>
      <li>
        WWN Primary Customer - maps WWN to a customer record to be used in standard export. Creating this rule for 1 record will cause all
        other records with the same WWN to go to override export. Tool provide reconciliation pop up to select the primary customer and will
        autogeneate this rule based on that selection.
      </li>
      <li>
        Ignore Loaded Host - regular expression with match group that causes the loaded hostname to be ignored when checking missmatch between
        loaded and decoded host. It is used to tell to export the record with decoded hostname despite the missmatch. Counter rule to tell to export 
        the WWN with the loaded hostname is host WWN rule that maps the WWN to loaded hostname. The tool allow one click setup of both these rules. 
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
