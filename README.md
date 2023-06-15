# STAMETS - STatic Analysis METricS

For programming language experts that rely on standard library static analyzers.
Requires Go 1.20

Currently supports gathering metrics for:
* Package loading:
    - R `Load` in `golang.org/x/tools/go/packages`
* SSA construction:
    - Replace all calls to `AllPackages`  in `golang.org/x/tools/go/ssautil` with `stamets.AllPackages`
* Standard Points-To Analysis (PTA).
    - Replace all calls to `Analyze` in `golang.org/x/tools/go/pointer` with `stamets.Analyze`
* Call graph metrics:
    - Call `GetCallGraphMetrics` by providing a `*callgraph.Graph` value e.g., produced by PTA

For wrappers around existing functions, the result is a metrics aggregator in the form of an appropriately
typed `Metrics` structure.
To extract the underlying result (and potential error), use the `Unpack` method.


## Infoermation about metrics

Gathered metrics include the following:
* Execution time
* **PTA**
    - Additional metrics are gathered for the sizes of points-to sets of the
queries included in the PTA results. These include: P50, P90, P99, Maximum size, Predominant points-to set size (mode)
* **Call graphs**
    - **Number of functions**
    - **Out-degree metrics**: P50, P90, P99, Maximum, Predominant out-degree (mode)
    - **In-degree metrics**: P50, P90, P99, Maximum, Predominant in-degree (mode)


Functions without out-going calls still contribute to out-degree metrics with a single 0 value.

## Example

Replacing a PTA `Analyze` call may be carried out as follows:
```
// Old:
// 1. Direct call to pointer.Analyze
ptaResults, err := pointer.Analyze(config)

// New:
// 1. Replace pointer.Analyze with stamets.Analyze
ptaMetrics := stamets.Analyze(config)
// 2. Unpack original results
ptaResults, err ptaMetrics := ptaMetrics.Unpack()
// 3. (Optional) Print metrics to stdout
fmt.Println(ptaMetrics)
```

The call graph produced by PTA may have its metrics extracted as follows:
```
pta, err := stamtes.Analyze(config).Unpack()
if err != nil {
    return
}

cgMetrics := stamets.GetCallGraphMetrics(pta.CallGraph)
```